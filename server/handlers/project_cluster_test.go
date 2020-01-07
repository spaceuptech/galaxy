package handlers

import (
	"log"
	"net/http"
	"testing"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/utils"
)

func TestHandleClusterRegistration(t *testing.T) {
	clusterRegistrationEndpoint := "http://localhost:4122/v1/galaxy/register-cluster"
	// TODO TOKEN CREATION IN FILE & SETTING IN HTTP HEADERS
	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    clusterRegistrationEndpoint,
	}
	tests := []struct {
		name          string
		httpBody      model.RegisterClusterPayload
		isErrExpected bool
	}{
		{
			name: "Register cluster test",
			httpBody: model.RegisterClusterPayload{
				ClusterID:  "cluster1",
				RunnerType: "runner1",
				Url:        "dummyUrl",
			},
			isErrExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Params = tt.httpBody
			err := utils.HttpRequest(h)
			if err != nil {
				t.Error(err)
			}
			log.Println("Response from space galaxy ", h.Response)
			// v := response.(map[string]interface{})
			// log.Println("v", v)
			// value, ok := v["error"]
			// if ok != tt.isErrExpected {
			// 	t.Errorf("error registering cluster with galaxy server %v", value)
			// }
		})
	}
}
