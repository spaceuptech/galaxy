package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
