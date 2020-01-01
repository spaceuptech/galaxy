package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

// TODO send TOKEN FROM HEADER

func TestHandleAddProject(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"
	createProjectEndpoint := "http://localhost:4050/v1/galaxy/project/create"

	testCases := []struct {
		name          string
		loginInfo     map[string]interface{}
		httpBody      *model.CreateProject
		isErrExpected bool
	}{
		{
			name: "Add project new1",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			httpBody: &model.CreateProject{
				ID: "new1",
				Environments: []*model.Environment{
					{
						ID:       "env1",
						Clusters: []*model.Cluster{&model.Cluster{ID: "cluster1"}},
					},
				},
			},
			isErrExpected: false,
		},
		{
			name: "Add project new2",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			httpBody: &model.CreateProject{
				ID: "new2",
				Environments: []*model.Environment{
					{
						ID:       "env2",
						Clusters: []*model.Cluster{{ID: "cluster2"}},
					},
				},
			},
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp, gotErr := utils.HttpRequest(http.MethodPost, loginEndpoint, nil, testCase.loginInfo, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				// t.Logf("token %s", token)
				_, gotErr := utils.HttpRequest(http.MethodPost, createProjectEndpoint, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}, testCase.httpBody, utils.SimpleRequest)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error add project got, %v wanted, %v", gotErr, testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleDeleteProject(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	testCases := []struct {
		name          string
		projectID     string
		loginInfo     map[string]interface{}
		isErrExpected bool
	}{
		{
			name:      "Open source login test",
			projectID: "new1",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			DeleteProjectEndpoint := fmt.Sprintf("http://localhost:4050/v1/galaxy/project/%s", testCase.projectID)

			resp, gotErr := utils.HttpRequest(http.MethodPost, loginEndpoint, nil, testCase.loginInfo, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				_, gotErr := utils.HttpRequest(http.MethodDelete, DeleteProjectEndpoint, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}, nil, utils.SimpleRequest)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error delete project got, %v wanted, %v", gotErr, testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleGetProject(t *testing.T) {
	// TODO REMOVE PRINT AND COMPARE TO THE ACTUAL VALUE IN TEST CASE
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"
	GetProjectEndpoint := "http://localhost:4050/v1/galaxy/project/new1"

	testCases := []struct {
		name          string
		loginInfo     map[string]interface{}
		isErrExpected bool
	}{
		{
			name: "Open source login test",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp, gotErr := utils.HttpRequest(http.MethodPost, loginEndpoint, nil, testCase.loginInfo, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				// t.Logf("token %s", token)
				resp, gotErr := utils.HttpRequest(http.MethodGet, GetProjectEndpoint, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}, nil, utils.SimpleRequest)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error delete project got, %v wanted, %v", gotErr, testCase.isErrExpected)
				}
				t.Log("Get project response", resp)
			}
		})
	}
}

func TestHandleGetProjects(t *testing.T) {
	// TODO REMOVE PRINT AND COMPARE TO THE ACTUAL VALUE IN TEST CASE
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"
	GetProjectsEndpoint := "http://localhost:4050/v1/galaxy/projects"

	testCases := []struct {
		name          string
		loginInfo     map[string]interface{}
		isErrExpected bool
	}{
		{
			name: "Open source login test",
			loginInfo: map[string]interface{}{
				"username": "admin",
				"key":      "1234",
			},
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp, gotErr := utils.HttpRequest(http.MethodPost, loginEndpoint, nil, testCase.loginInfo, utils.SimpleRequest)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				// t.Logf("token %s", token)
				resp, gotErr := utils.HttpRequest(http.MethodGet, GetProjectsEndpoint, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}, nil, utils.SimpleRequest)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error delete project got, %v wanted, %v", gotErr, testCase.isErrExpected)
				}
				t.Log("Get project response", resp)
			}
		})
	}
}
