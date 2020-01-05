package cmd

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/galaxy/model"
)

//CodeStart starts the code commands
func CodeStart(envID string) (*model.Service, error) {
	credential, err := getCreds()
	if err != nil {
		return nil, err
	}

	selectedAccount := getSelectedAccount(credential)

	loginRes, err := login(selectedAccount)
	if err != nil {
		return nil, err
	}

	c, err := getServiceConfig(".galaxy.yaml")
	if err != nil {
		c, err = generateServiceConfig(loginRes.Projects, selectedAccount, envID)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func generateServiceConfig(projects []model.Projects, selectedaccount *model.Account, envID string) (*model.Service, error) {
	progLang, err := getProgLang()
	if err != nil {
		return nil, err
	}
	var envNameID string
	var clusters []string
	serviceName := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Service Name"}, &serviceName); err != nil {
		return nil, err
	}
	defaultServiceID := strings.ReplaceAll(serviceName, " ", "-")
	serviceID := ""
	if err := survey.AskOne(&survey.Input{Message: "Enter Service ID", Default: strings.ToLower(defaultServiceID)}, &serviceID); err != nil {
		return nil, err
	}
	var port int32
	if err := survey.AskOne(&survey.Input{Message: "Enter Service Port"}, &port); err != nil {
		return nil, err
	}
	projectNameID := ""
	if err := survey.AskOne(&survey.Select{Message: "Select Project Name", Options: getProjects(projects)}, &projectNameID); err != nil {
		return nil, err
	}

	temp := strings.Split(projectNameID, " ")
	projectID := temp[0]

	var project model.Projects
	if envID == "none" {
		project, err := getProject(projectID, projects)
		if err != nil {
			return nil, err
		}
		if err := survey.AskOne(&survey.Select{Message: "Select Environment", Options: getEnvironments(project)}, &envNameID); err != nil {
			return nil, err
		}
		temp := strings.Split(envNameID, " ")
		envID = temp[0]
	}

	selectedEnv, err := getEnvironment(envID, project.Environments)
	if err != nil {
		return nil, err
	}
	if err := survey.AskOne(&survey.MultiSelect{Message: "Select Clusters", Options: getClusters(selectedEnv)}, &clusters); err != nil {
		return nil, err
	}
	var progCmd string
	if err := survey.AskOne(&survey.Input{Message: "Enter Run Cmd",
		Default: strings.Join(getCmd(progLang), " ")}, &progCmd); err != nil {
		return nil, err
	}
	progCmds := strings.Split(progCmd, " ")
	img, err := getImage(progLang)
	if err != nil {
		return nil, err
	}

	c := &model.Service{
		ID:          serviceID,
		Name:        serviceName,
		ProjectID:   projectID,
		Environment: envID,
		Version:     "v1",
		Scale:       model.ScaleConfig{Replicas: 0, MinReplicas: 0, MaxReplicas: 100, Concurrency: 50},
		Tasks: []model.Task{
			{
				ID:        serviceID,
				Name:      serviceName,
				Ports:     []model.Port{model.Port{Protocol: "http", Port: port}},
				Resources: model.Resources{CPU: 250, Memory: 512},
				Docker:    model.Docker{Image: img, Cmd: progCmds},
				Env:       map[string]string{"URL": selectedaccount.ServerUrl},
			},
		},
		Whitelist: []string{"project:*"},
		Upstreams: []model.Upstream{model.Upstream{ProjectID: projectID, Service: "*"}},
		Runtime:   "code",
	}
	return c, nil
}
