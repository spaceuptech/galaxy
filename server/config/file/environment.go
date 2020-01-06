package file

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/launchpad/model"
)

// SetDefaultEnvironment is used to set default environment of specified project
func (m *Manager) SetDefaultEnvironment(ctx context.Context, accountID, projectID, defaultEnv string) error {
	// get specified project info from projects table
	projects, err := m.GetProject(ctx, projectID, accountID)
	if err != nil {
		logrus.Errorf("error setting default environment - %v", err)
		return err
	}

	// there should only be a single project for specified projectID
	if len(projects) == 1 {
		projects[0].DefaultEnv = defaultEnv
		return m.updateProject(accountID, &projects[0])
	}
	logrus.Errorf("error setting default environment project length not equal to one")
	return err
}

// AddEnvironment is used add specified environment to the project config if it doesn't exits
func (m *Manager) AddEnvironment(ctx context.Context, accountID, projectID string, req *model.Environment) error {
	// get specified project from database
	projects, err := m.GetProject(ctx, projectID, accountID)
	if err != nil {
		logrus.Errorf("error adding environment - %v", err)
		return err
	}

	// there should only be a single project for specified projectID
	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		// convert to envs as it is stored as json string in projects table
		if err := json.Unmarshal([]byte(project.Environments), &envs); err != nil {
			logrus.Errorf("error adding environment unable to unmarshal envs - %v", err)
			return err
		}

		isEnvFound := false
		// check if environment already exists
		for _, environment := range envs {
			isEnvFound = true
			if environment.ID == req.ID {
				logrus.Errorf("error adding environment specified environment already exists")
				return err
			}
		}

		// if doesn't exits then add new environment & update the database
		if isEnvFound {
			envs = append(envs, req)
			// convert envs to json string as it stored as json string in projects table
			data, err := json.Marshal(envs)
			if err != nil {
				logrus.Errorf("error adding environment unable to marshal envs - %v", err)
				return err
			}
			project.Environments = string(data)
			return m.updateProject(accountID, &project)
		}
	}
	logrus.Errorf("error adding environment project length not equal to one")
	return err
}

// DeleteEnvironment deletes specified environment from project config if it exits
func (m *Manager) DeleteEnvironment(ctx context.Context, accountID, projectID, environmentID string) error {
	// get specified project from database
	projects, err := m.GetProject(ctx, projectID, accountID)
	if err != nil {
		logrus.Errorf("error deleting environment - %v", err)
		return err
	}

	// there should only be a single project for specified projectID
	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		// convert to envs as it is stored as json string in projects table
		if err := json.Unmarshal([]byte(project.Environments), &envs); err != nil {
			logrus.Errorf("error deleting environment unable to unmarshal envs - %v", err)
			return err
		}

		// check if environment already exists
		isEnvFound := false
		for environmentIndex, environment := range envs {
			if environment.ID == environmentID {
				isEnvFound = true
				envs = removeEnvironmentAtIndex(envs, environmentIndex)
				// convert envs to json string before updating table
				data, err := json.Marshal(envs)
				if err != nil {
					logrus.Errorf("error deleting environment unable to marshal envs - %v", err)
					return err
				}
				project.Environments = string(data)
				return m.updateProject(accountID, &project)
			}
		}

		// throw error if specified environment doesn't exists
		if !isEnvFound {
			logrus.Errorf("error deleting environment specified environment not found")
			return err
		}
	}
	logrus.Errorf("error deleting environment specified project does not exist")
	return err
}
