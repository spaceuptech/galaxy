package runner

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

func (runner *Runner) handleCreateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// Verify token
		_, err := runner.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}
		// Parse request body
		project := new(model.Project)
		if err := json.NewDecoder(r.Body).Decode(project); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// Apply the service config
		if err := runner.driver.CreateProject(project); err != nil {
			logrus.Errorf("Failed to create project - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}

func (runner *Runner) handleServiceRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// Verify token
		_, err := runner.auth.VerifyToken(utils.GetToken(r))
		if err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}
		// Parse request body
		service := new(model.Service)
		if err := json.NewDecoder(r.Body).Decode(service); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		// TODO: Override the project id present in the service object with the one present in the token if user not admin

		// Apply the service config
		if err := runner.driver.ApplyService(service); err != nil {
			logrus.Errorf("Failed to apply service - %s", err.Error())
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		utils.SendEmptySuccessResponse(w, r)
	}
}

func (runner *Runner) handleProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		// http: Request.RequestURI can't be set in client requests.
		// http://golang.org/src/pkg/net/http/client.go
		r.RequestURI = ""

		// Change the destination with the original host and port
		ogHost := r.Header.Get("x-og-host")
		ogPort := r.Header.Get("x-og-port")
		r.Host = ogHost
		r.URL.Host = fmt.Sprintf("%s:%s", ogHost, ogPort)

		// Retrieve project id and service
		array := strings.Split(ogHost, ".")
		project, service := array[1], array[0]

		// Set the url scheme to http
		r.URL.Scheme = "http"

		// Add to active request count
		// TODO: add support for multiple versions
		runner.chAppend <- &model.ProxyMessage{Service: service, Project: project, Version: "v1", NodeID: "runner-proxy", ActiveRequests: 1}

		// Wait for the service to scale up
		if err := runner.debounce.Wait(fmt.Sprintf("proxy-%s-%s", project, service), func() error {
			return runner.driver.WaitForService(project, service)
		}); err != nil {
			utils.SendErrorResponse(w, r, http.StatusServiceUnavailable, err)
			return
		}
		var res *http.Response
		for i := 0; i < 5; i++ {
			// Fire the request
			var err error
			res, err = http.DefaultClient.Do(r)
			if err != nil {
				utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
				return
			}
			// TODO: Make this retry logic better
			if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusServiceUnavailable {
				break
			}
			time.Sleep(350 * time.Millisecond)

			// Close the body
			_, _ = io.Copy(ioutil.Discard, res.Body)
			utils.CloseReaderCloser(res.Body)
		}

		defer utils.CloseReaderCloser(res.Body)
		// Copy headers and status code
		w.WriteHeader(res.StatusCode)
		for k, v := range res.Header {
			w.Header().Set(k, v[0])
		}
		_, _ = io.Copy(w, res.Body)
	}
}

func (runner *Runner) handleDatabaseService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Close the body of the request
		defer utils.CloseReaderCloser(r.Body)

		service := new(model.ManagedService)
		_ = json.NewDecoder(r.Body).Decode(service)
	}
}
