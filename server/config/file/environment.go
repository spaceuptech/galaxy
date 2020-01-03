package file

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spaceuptech/launchpad/model"
)

// AddEnvironment adds new environment to the specified project
func (m *Manager) AddEnvironment(ctx context.Context, projectID string, req *model.Environment) error {
	// get specified project from database
	projects, err := m.GetProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("error adding environment - %v", err)
	}

	// there should only be a single project for specified projectID
	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		// unmarshal environment as it is stored as json string in database
		if err := json.Unmarshal([]byte(project.Environments), envs); err != nil {
			return fmt.Errorf("error adding environment unable to unmarshal envs - %v", err)
		}

		isEnvFound := false
		// check if environment already exists
		for _, environment := range envs {
			isEnvFound = true
			if environment.ID == req.ID {
				fmt.Errorf("error adding environment specified environment already exists")
			}
		}

		// if doesn't exits then add new environment & update the database
		if isEnvFound {
			envs = append(envs, req)
			data, err := json.Marshal(envs)
			if err != nil {
				return fmt.Errorf("error adding environment unable to marshal envs - %v", err)
			}
			project.Environments = string(data)
			return m.updateProject(project)
		}
	}
	return fmt.Errorf("error adding environment project length not equal to one")
}

// DeleteEnvironment deletes specified environment from database if it exists
func (m *Manager) DeleteEnvironment(ctx context.Context, projectID, environmentID string) error {
	projects, err := m.GetProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("error deleting environment - %v", err)
	}

	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		if err := json.Unmarshal([]byte(project.Environments), envs); err != nil {
			return fmt.Errorf("error deleting environment unable to unmarshal envs - %v", err)
		}

		isEnvFound := false
		for _, environment := range envs {
			if environment.ID == environmentID {
				isEnvFound = true
				// TODO REMOVE ENVIRONMENT HERE
				data, err := json.Marshal(envs)
				if err != nil {
					return fmt.Errorf("error deleting environment unable to marshal envs - %v", err)
				}
				project.Environments = string(data)
				return m.updateProject(project)
			}
		}
		if !isEnvFound {
			return fmt.Errorf("error deleting environment specified environment not found")
		}
	}
	return fmt.Errorf("error deleting environment project length not equal to one")
}
