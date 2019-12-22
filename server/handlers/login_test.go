package handlers

import (
	"testing"

	"github.com/spaceuptech/launchpad/utils"
)

func TestHandleLogin(t *testing.T) {
	loginEndpoint := "http://localhost:4050/v1/galaxy/login"

	testCases := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Open source login test",
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, gotErr := utils.HttpRequest("", loginEndpoint, utils.SimpleRequest)
			if (gotErr != nil) != testCase.wantErr {
				t.Errorf("Error login got, %v wanted, %v", gotErr, testCase.wantErr)
			}
		})
	}
}
