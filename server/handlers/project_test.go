package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/utils"
)

func TestHandleAddProject(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"
	createProjectEndpoint := "http://localhost:4050/v1/galaxy/project/create"

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
		httpBody      *model.CreateProject
		isErrExpected bool
	}{
		{
			name: "Add project new1",
			httpBody: &model.CreateProject{
				ID:                 "new6",
				DefaultEnvironment: "production",
				Environments: []*model.Environment{
					{
						ID:       "env1",
						Clusters: []*model.Cluster{{ID: "cluster1"}},
					},
				},
			},
			isErrExpected: false,
		},
		// {
		// 	name: "Add project new2",
		// 	httpBody: &model.CreateProject{
		// 		ID: "new2",
		// 		Environments: []*model.Environment{
		// 			{
		// 				ID:       "env2",
		// 				Clusters: []*model.Cluster{{ID: "cluster2"}},
		// 			},
		// 		},
		// 	},
		// 	isErrExpected: false,
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := map[string]interface{}{}
			h.Response = &resp
			gotErr := utils.HttpRequest(h)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v wanted - %v, %v", gotErr, resp["error"], testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				h.Url = createProjectEndpoint
				h.Params = testCase.httpBody
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				gotErr := utils.HttpRequest(h)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error add project, %v wanted - %v, %v", gotErr, resp["error"], testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleDeleteProject(t *testing.T) {
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
		isErrExpected bool
	}{
		{
			name:          "Open source login test",
			projectID:     "new1",
			isErrExpected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resp := map[string]interface{}{}
			h.Response = &resp
			gotErr := utils.HttpRequest(h)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("Error login got, %v  wanted, %v", gotErr, testCase.isErrExpected)
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				h.Method = http.MethodDelete
				h.Url = fmt.Sprintf("http://localhost:4050/v1/galaxy/project/%s", testCase.projectID)
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				gotErr := utils.HttpRequest(h)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error delete project, %v wanted - %v, %v", gotErr, resp["error"], testCase.isErrExpected)
				}
			}
		})
	}
}

func TestHandleGetProject(t *testing.T) {
	// TODO REMOVE PRINT AND COMPARE TO THE ACTUAL VALUE IN TEST CASE
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
		projectId     string
		isErrExpected bool
	}{
		{
			name:          "Open source login test",
			projectId:     "new2",
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
				return
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				h.Url = fmt.Sprintf("http://localhost:4050/v1/galaxy/project/%s", testCase.projectId)
				h.Method = http.MethodGet
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				gotErr := utils.HttpRequest(h)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error getting project, %v wanted - %v, %v", gotErr, resp["error"], testCase.isErrExpected)
					return
				}
				t.Log("Get project response", h.Response)
			}
		})
	}
}

func TestHandleGetProjects(t *testing.T) {
	// TODO REMOVE PRINT AND COMPARE TO THE ACTUAL VALUE IN TEST CASE
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"
	GetProjectsEndpoint := "http://localhost:4050/v1/galaxy/projects"

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
		isErrExpected bool
	}{
		{
			name:          "Open source login test",
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
				return
			}

			if token, ok := resp["token"]; !ok {
				t.Logf("token not found")
			} else {
				h.Url = GetProjectsEndpoint
				h.Method = http.MethodGet
				h.Headers = map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
				gotErr := utils.HttpRequest(h)
				if (gotErr != nil) != testCase.isErrExpected {
					t.Errorf("Error getting projects, %v wanted - %v, %v", gotErr, resp["error"], testCase.isErrExpected)
					return
				}
				t.Log("Get project response", h.Response)
			}
		})
	}
}
