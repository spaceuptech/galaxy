package handlers

import (
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp, gotErr := utils.HttpRequest(testCase.httpBody, loginEndpoint, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}
			token, ok := resp.(string)
			if !ok {
				t.Logf("Generated Token - %s", token)
			}
		})
	}
}
