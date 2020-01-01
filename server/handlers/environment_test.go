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
			projectID:     "new1",
			environmentID: "evn",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			httpBody: &model.Environment{
				ID: "env2",
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
			DeleteProjectEndpoint := fmt.Sprintf("http://localhost:4050/v1/galaxy/project/%s/%s", testCase.projectID, testCase.environmentID)

			resp, gotErr := utils.HttpRequest(http.MethodPost, loginEndpoint, nil, testCase.loginInfo, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				// t.Logf("token %s", token)
				_, gotErr := utils.HttpRequest(http.MethodPost, DeleteProjectEndpoint, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}, testCase.httpBody, utils.SimpleRequest)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error delete project got, %v wanted, %v", gotErr, testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleDeleteEnvironment(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	testCases := []struct {
		name          string
		projectID     string
		environmentID string
		loginInfo     map[string]interface{}
		isErrExpected bool
	}{
		{
			name:          "Deleting environment env1",
			projectID:     "new1",
			environmentID: "env2",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ProjectEnvironmentEndpoint := fmt.Sprintf("http://localhost:4050/v1/galaxy/project/%s/%s", testCase.projectID, testCase.environmentID)

			resp, gotErr := utils.HttpRequest(http.MethodPost, loginEndpoint, nil, testCase.loginInfo, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				_, gotErr := utils.HttpRequest(http.MethodDelete, ProjectEnvironmentEndpoint, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}, nil, utils.SimpleRequest)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error delete project got, %v wanted, %v", gotErr, testCase.isErrExpected)
				}
			}
		})
	}
}
