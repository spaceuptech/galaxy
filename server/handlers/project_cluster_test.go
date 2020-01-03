package handlers

import (
	"log"
	"net/http"
	"testing"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

func TestHandleClusterRegistration(t *testing.T) {
	clusterRegistrationEndpoint := "http://localhost:4122/v1/galaxy/register-cluster"
	// TODO TOKEN CREATION IN FILE & SETTING IN HTTP HEADERS
	tests := []struct {
		name          string
		httpBody      model.RegisterClusterRequest
		isErrExpected bool
	}{
		{
			name: "Register cluster test",
			httpBody: model.RegisterClusterRequest{
				ClusterID:  "cluster1",
				RunnerType: "runner1",
				Url:        "dummyUrl",
			},
			isErrExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := utils.HttpRequest(http.MethodPost, clusterRegistrationEndpoint, nil, tt.httpBody, utils.SimpleRequest)
			if err != nil {
				t.Error(err)
			}
			log.Println("Response from space galaxy ", response)
			// v := response.(map[string]interface{})
			// log.Println("v", v)
			// value, ok := v["error"]
			// if ok != tt.isErrExpected {
			// 	t.Errorf("error registering cluster with galaxy server %v", value)
			// }
		})
	}
}
