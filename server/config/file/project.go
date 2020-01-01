package file

import (
	"context"
	"fmt"

	"github.com/spaceuptech/launchpad/model"
)

// AddProject adds a project to config file if it doesn't exist
func (m *Manager) AddProject(ctx context.Context, req *model.CreateProject) error {
	m.RLock()
	config := m.galaxyConfig
	m.RUnlock()

	for _, project := range config.Projects {
		if project.ID == req.ID {
			return fmt.Errorf("error adding project, project already exists in galaxy config")
		}
	}

	config.Projects = append(config.Projects, req)

	return m.StoreConfigToFile(ctx, config)
}

// GetProject returns the specified project if it exists in config file
func (m *Manager) GetProject(ctx context.Context, projectID string) (*model.CreateProject, error) {
	m.Lock()
	defer m.Unlock()

	config := m.galaxyConfig

	for _, project := range config.Projects {
		if project.ID == projectID {
			return project, nil
		}
	}

	return nil, fmt.Errorf("error getting project project doesn't exist in config")
}

// GetProjects return all the projects present in config file
func (m *Manager) GetProjects(ctx context.Context) ([]*model.CreateProject, error) {
	m.Lock()
	defer m.Unlock()

	if m.galaxyConfig != nil {
		return m.galaxyConfig.Projects, nil
	}

	return nil, fmt.Errorf("error getting projects galaxy config not initialzed")
}

// DeleteProject deletes a specified project from the config file
func (m *Manager) DeleteProject(ctx context.Context, projectID string) error {
	// TODO WILL THERE BE RACE CONDITION
	m.Lock()
	config := m.galaxyConfig
	m.Unlock()

	for i, project := range config.Projects {
		if project.ID == projectID {
			newProject, err := remove(config.Projects, i)
			if err != nil {
				return err
			}
			config.Projects = newProject.([]*model.CreateProject)
			return m.StoreConfigToFile(ctx, config)
		}
	}

	return fmt.Errorf("error deleting project, project doesn't exist in config")
}
