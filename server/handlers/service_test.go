package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/server/config"
	"github.com/spaceuptech/launchpad/utils"
	"github.com/spaceuptech/launchpad/utils/auth"
)

func TestHandleApplyService(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    loginEndpoint,
		Params: map[string]interface{}{
			"username": "admin",
			"key":      "1234",
		},
	}

	testCases := []struct {
		name          string
		serviceID     string
		httpBody      *model.Service
		isErrExpected bool
	}{
		{
			name:          "Adding environment env2",
			serviceID:     "new2",
			httpBody:      &model.Service{},
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := map[string]interface{}{}
			h.Response = &resp
			gotErr := utils.HttpRequest(h)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				h.Url = fmt.Sprintf("http://localhost:4050/v1/galaxy/service/%s/ui", testCase.serviceID)
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				h.Params = testCase.httpBody
				if gotErr := utils.HttpRequest(h); (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error adding environment got, %v -%v wanted, %v", gotErr, resp["error"], testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleDeleteService(t *testing.T) {
	type args struct {
		auth         *auth.Module
		galaxyConfig *config.Module
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleDeleteService(tt.args.auth, tt.args.galaxyConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleDeleteService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleUpsertCLiService(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    loginEndpoint,
		Params: map[string]interface{}{
			"username": "admin",
			"key":      "1234",
		},
	}

	testCases := []struct {
		name          string
		serviceID     string
		httpBody      *model.FileStoreEventPayload
		isErrExpected bool
	}{
		{
			name:          "testing cli upsert service",
			serviceID:     "service1",
			isErrExpected: false,
			httpBody: &model.FileStoreEventPayload{
				Data: &model.FileStoreData{
					Path: "",
					Meta: &model.ServiceRequest{
						IsDeploy: true,
						Service: &model.Service{
							Clusters:    []string{"india", "usa", "uk"},
							Environment: "production",
							ID:          "service1",
							ProjectID:   "new2",
							Version:     "0.0.1",
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := map[string]interface{}{}
			h.Response = &resp
			gotErr := utils.HttpRequest(h)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				h.Url = fmt.Sprintf("http://localhost:4050/v1/galaxy/service/%s/cli", testCase.serviceID)
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				h.Params = testCase.httpBody
				if gotErr := utils.HttpRequest(h); (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error upserting cli service got, %v -%v wanted, %v", gotErr, resp["error"], testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleUpsertUIService(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    loginEndpoint,
		Params: map[string]interface{}{
			"username": "admin",
			"key":      "1234",
		},
	}

	testCases := []struct {
		name          string
		serviceID     string
		httpBody      *model.Service
		isErrExpected bool
	}{
		{
			name:          "testing cli upsert service",
			serviceID:     "new2",
			httpBody:      &model.Service{},
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := map[string]interface{}{}
			h.Response = &resp
			gotErr := utils.HttpRequest(h)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				h.Url = fmt.Sprintf("http://localhost:4050/v1/galaxy/service/ui%s", testCase.serviceID)
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				h.Params = testCase.httpBody
				if gotErr := utils.HttpRequest(h); (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error upserting ui servicet got, %v -%v wanted, %v", gotErr, resp["error"], testCase.isErrExpected)
				}
			}
		})
	}
}
