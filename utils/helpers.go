package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spaceuptech/launchpad/model"
)

// FireGraphqlQuery is a general function used to send graphql queries to space cloud
func FireGraphqlQuery(params *model.InsertRequest, responseType int) (interface{}, error) {
	// TODO: remove graphql endpoint field from billing moudle & make it a constant
	requestBody := new(bytes.Buffer)
	if err := json.NewEncoder(requestBody).Encode(params); err != nil {
		return nil, err
	}

	resp, err := http.Post(GraphqlEndpoint, ApplicationJson, requestBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch responseType {
	case GraphqlMutation:
		v := model.MutationQueryResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return nil, err
		}

		if resp.StatusCode != 200 && (v.Data.Insert.Status != 200 || v.Data.Update.Status != 200) {
			return nil, fmt.Errorf("error while inserting data in database")
		}
		return v, nil
	default:
		return nil, fmt.Errorf("invalid response type")
	}
}

type HttpModel struct {
	Method   string
	Url      string
	Headers  map[string]string
	Params   interface{}
	Response interface{}
}

// HttpRequest is a general function for sending http request to the provided url
func HttpRequest(h *HttpModel) error {
	if h.Response == nil {
		h.Response = &map[string]interface{}{}
	}
	// encode to json
	requestBody := new(bytes.Buffer)
	if h.Params != nil {
		if err := json.NewEncoder(requestBody).Encode(h.Params); err != nil {
			return fmt.Errorf("error encoding body for http request %v", err)
		}
	}

	client := http.Client{}

	// create http request
	httpRequest, err := http.NewRequest(h.Method, h.Url, requestBody)
	if err != nil {
		return fmt.Errorf("error creating http request - %v", err)
	}

	// set http headers
	if h.Headers != nil {
		for key, value := range h.Headers {
			httpRequest.Header.Add(key, value)
		}
	}

	// make http request
	resp, err := client.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("error sending http %s request - %v", h.Method, err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(h.Response)
	if resp.StatusCode != http.StatusOK {
		log.Println("error", h.Response)
		return fmt.Errorf("error response from server with status code %v", resp.StatusCode)
	}
	return nil

}

// GetTokenFromHeader returns the token from the request header
func GetTokenFromHeader(r *http.Request) string {
	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}
	return strings.TrimPrefix(tokens[0], "Bearer ")
}
