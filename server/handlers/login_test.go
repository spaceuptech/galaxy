package handlers

import (
	"net/http"
	"testing"

	"github.com/spaceuptech/launchpad/utils"
)

func TestHandleLogin(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	testCases := []struct {
		name          string
		httpBody      map[string]interface{}
		isErrExpected bool
	}{
		{
			name: "Open source login test",
			httpBody: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			isErrExpected: false,
		},
		{
			name: "Invalid credentials",
			httpBody: map[string]interface{}{
				"username": "nil",
				"key":      "nil",
			},
			isErrExpected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp, gotErr := utils.HttpRequest(http.MethodPost, loginEndpoint, nil, testCase.httpBody, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if !testCase.isErrExpected {
				if _, ok := resp["token"]; !ok {
					t.Logf("token not found")
				}
			}

		})
	}
}
