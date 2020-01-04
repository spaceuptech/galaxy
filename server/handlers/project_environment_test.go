package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

func TestHandleAddEnvironment(t *testing.T) {
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
		projectID     string
		environmentID string
		loginInfo     map[string]interface{}
		httpBody      *model.Environment
		isErrExpected bool
	}{
		{
			name:          "Adding environment env2",
			projectID:     "new2",
			environmentID: "evn12",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			httpBody: &model.Environment{
				ID: "env12",
				Clusters: []*model.Cluster{
					{
						ID: "cluster1",
					},
				},
			},
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
				h.Url = fmt.Sprintf("http://localhost:4050/v1/galaxy/project/%s/%s", testCase.projectID, testCase.environmentID)
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				h.Params = testCase.httpBody
				gotErr := utils.HttpRequest(h)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error adding environment got, %v -%v wanted, %v", gotErr, resp["error"], testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleDeleteEnvironment(t *testing.T) {
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
		projectID     string
		environmentID string
		isErrExpected bool
	}{
		{
			name:          "Deleting environment env1",
			projectID:     "new2",
			environmentID: "env12",
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
				h.Method = http.MethodDelete
				h.Url = fmt.Sprintf("http://localhost:4050/v1/galaxy/project/%s/%s", testCase.projectID, testCase.environmentID)
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				gotErr := utils.HttpRequest(h)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error deleting environment got, %v -%v wanted, %v", gotErr, resp["error"], testCase.isErrExpected)
				}
			}
		})
	}
}
