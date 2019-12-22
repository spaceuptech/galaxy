package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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

// HttpRequest is a general function for sending http request to the provided url
func HttpRequest(params interface{}, url string, functionCallType int) (interface{}, error) {
	requestBody := new(bytes.Buffer)
	if err := json.NewEncoder(requestBody).Encode(params); err != nil {
		return nil, fmt.Errorf("error encoding body for http request %v", err)
	}

	resp, err := http.Post(url, ApplicationJson, requestBody)
	if err != nil {
		return nil, fmt.Errorf("error sending http request for pinging the cluster %v", err)
	}
	defer resp.Body.Close()

	switch functionCallType {
	case Ping:
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error pinging cluster")
		}
		return nil, nil

	case SimpleRequest:
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error response from server with status code %v", resp.StatusCode)
		}
		body := new(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(body)
		return body, nil

	default:
		return nil, fmt.Errorf("invalid response type")
	}
}
