package handlers

import (
	"net/http"
	"testing"

	"github.com/spaceuptech/launchpad/utils"
)

func TestHandleLogin(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    loginEndpoint,
	}

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
			h.Params = testCase.httpBody
			resp := map[string]interface{}{}
			h.Response = &resp
			gotErr := utils.HttpRequest(h)
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
