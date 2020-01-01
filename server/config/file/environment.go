package file

import (
	"context"
	"fmt"

	"github.com/spaceuptech/launchpad/model"
)

func (m *Manager) AddEnvironment(ctx context.Context, projectID string, req *model.Environment) error {
	m.RLock()
	config := m.galaxyConfig
	m.RUnlock()

	isProjectFound := false
	projectIndex := 0

	for p, project := range config.Projects {
		if project.ID == projectID {
			projectIndex = p
			isProjectFound = true
			for _, environment := range project.Environments {
				if environment.ID == req.ID {
					return fmt.Errorf("error adding environment, specified environment already exists in galaxy config")
				}
			}
		}
	}

	// add new environment
	if isProjectFound {
		config.Projects[projectIndex].Environments = append(config.Projects[projectIndex].Environments, req)
		return m.StoreConfigToFile(ctx, config)
	}

	return fmt.Errorf("error adding environment, project doesn't exists in galaxy config")
}

// func (m *Manager) SetEnvironment(ctx context.Context, projectID string, req *model.Environment) error {
// 	m.RLock()
// 	config := m.galaxyConfig
// 	m.RUnlock()
//
// 	isEnvFound := false
//
// 	for _, project := range config.Projects {
// 		if project.ID == projectID {
// 			isEnvFound = true
// 			for i, environment := range project.Environments {
// 				if environment.ID == req.ID {
// 					// update
// 					project.Environments[i] = req
// 					break
// 				}
// 			}
// 		}
// 	}
//
// 	if !isEnvFound {
// 		return fmt.Errorf("error adding environment, project doesn't exists in galaxy config")
// 	}
//
// 	return m.StoreConfigToFile(ctx, config)
// }

func (m *Manager) DeleteEnvironment(ctx context.Context, projectID, environmentId string) error {
	m.RLock()
	config := m.galaxyConfig
	m.RUnlock()

	isEnvFound := false

	for _, project := range config.Projects {
		if project.ID == projectID {
			for i, environment := range project.Environments {
				if environment.ID == environmentId {
					isEnvFound = true
					newEnvironment, err := remove(project.Environments, i)
					if err != nil {
						return err
					}
					project.Environments = newEnvironment.([]*model.Environment)
					return m.StoreConfigToFile(ctx, config)
				}
			}
		}
	}

	if !isEnvFound {
		return fmt.Errorf("error removing environment, specified environment doesn't exists in galaxy config")
	}
	return nil
}
