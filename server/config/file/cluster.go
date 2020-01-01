package file

import (
	"context"
	"fmt"

	"github.com/spaceuptech/launchpad/model"
)

func (m *Manager) AddCluster(ctx context.Context, projectID, environmentID string, req *model.Cluster) error {
	m.RLock()
	config := m.galaxyConfig
	m.RUnlock()

	isEnvFound := false
	projectIndex, environmentIndex := 0, 0

	for p, project := range config.Projects {
		if project.ID == projectID {
			projectIndex = p
			for e, environment := range project.Environments {
				if environment.ID == environmentID {
					environmentIndex = e
					isEnvFound = true
					for _, cluster := range environment.Clusters {
						if cluster.ID == req.ID {
							return fmt.Errorf("error adding cluster, provided cluster already exists in galaxy config")
						}
					}
				}
			}
		}
	}

	if isEnvFound {
		newCluster := new(model.Cluster)
		config.Projects[projectIndex].Environments[environmentIndex].Clusters = append(config.Projects[projectIndex].Environments[environmentIndex].Clusters, newCluster)
		return m.StoreConfigToFile(ctx, config)
	}

	return fmt.Errorf("error adding environment, specified environment doesn't exists in galaxy config")
}

func (m *Manager) DeleteCluster(ctx context.Context, projectID, environmentID, clusterID string) error {
	m.RLock()
	config := m.galaxyConfig
	m.RUnlock()

	isEnvFound := false

	for _, project := range config.Projects {
		if project.ID == projectID {
			for _, environment := range project.Environments {
				if environment.ID == environmentID {
					isEnvFound = true
					for i, cluster := range environment.Clusters {
						if cluster.ID == clusterID {
							newClusters, err := remove(environment.Clusters, i)
							if err != nil {
								return err
							}
							environment.Clusters = newClusters.([]*model.Cluster)
							return m.StoreConfigToFile(ctx, config)
						}
					}
				}
			}
		}
	}

	if !isEnvFound {
		return fmt.Errorf("error removing cluster, specified environment doesn't exists in galaxy config")
	}
	return nil
}
