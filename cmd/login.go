package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/utils"
)

func login(selectedAccount *model.Account) (*model.LoginResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"user": selectedAccount.UserName,
		"key":  selectedAccount.Key,
	})
	if err != nil {
		logrus.Errorf("error in login unable to marshal data for login got error message - %v", err)
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/v1/config/login", selectedAccount.ServerUrl), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logrus.Errorf("error in login unable to send http request for login got error message - %v", err)
		return nil, err
	}
	defer utils.CloseReaderCloser(resp.Body)

	loginResp := new(model.LoginResponse)
	_ = json.NewDecoder(resp.Body).Decode(loginResp)

	if resp.StatusCode != 200 {
		logrus.Errorf("error in login unable to login got status code %v with error message - %v", resp.StatusCode, loginResp.Error)
		return nil, fmt.Errorf("%v", loginResp.Error)
	}
	return loginResp, nil
}

// LoginStart logs the user in galaxy
func LoginStart(userName, key, url string, local bool) error {
	if userName == "None" {
		if err := survey.AskOne(&survey.Input{Message: "Enter username:"}, &userName); err != nil {
			logrus.Error("error starting login unable to get username from user got error message -%v", err)
			return err
		}
	}
	if key == "None" {
		if err := survey.AskOne(&survey.Password{Message: "Enter key:"}, &key); err != nil {
			logrus.Error("error starting login unable to get key from user got error message -%v", err)
			return err
		}
	}
	selectedAccount := model.Account{
		UserName:  userName,
		Key:       key,
		ServerUrl: url,
	}
	loginRes, err := login(&selectedAccount)
	if err != nil {
		return err
	}
	fmt.Printf("Login Successful\n")
	selectedAccount = model.Account{
		ID:        loginRes.Token,
		UserName:  userName,
		Key:       key,
		ServerUrl: url,
	}
	if err := checkCred(&selectedAccount); err != nil {
		return err
	}
	return nil
}
