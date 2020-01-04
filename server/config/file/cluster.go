package file

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

func (m *Manager) AddProjectCluster(ctx context.Context, accountID, projectID, environmentID string, req *model.Cluster) error {
	projects, err := m.GetProject(ctx, projectID, accountID)
	if err != nil {
		return fmt.Errorf("error adding cluster - %v", err)
	}

	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		if err := json.Unmarshal([]byte(project.Environments), &envs); err != nil {
			return fmt.Errorf("error adding cluster unable to unmarshal envs - %v", err)
		}

		isClusterFound, environmentIndex := false, 0
		for _, environment := range envs {
			if environment.ID == environmentID {
				for _, cluster := range environment.Clusters {
					isClusterFound = true
					if cluster.ID == req.ID {
						fmt.Errorf("error adding cluster specified cluster already exists")
					}
				}
			}
		}
		if isClusterFound {
			envs[environmentIndex].Clusters = append(envs[environmentIndex].Clusters, req)
			data, err := json.Marshal(envs)
			if err != nil {
				return fmt.Errorf("error adding cluster unable to marshal envs - %v", err)
			}
			project.Environments = string(data)
			return m.updateProject(accountID, &project)
		}
	}
	return fmt.Errorf("error adding cluster project length not equal to one")
}

func (m *Manager) DeleteProjectCluster(ctx context.Context, accountID, projectID, environmentID, clusterID string) error {
	projects, err := m.GetProject(ctx, projectID, accountID)
	if err != nil {
		return fmt.Errorf("error adding cluster - %v", err)
	}

	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		if err := json.Unmarshal([]byte(project.Environments), &envs); err != nil {
			return fmt.Errorf("error adding cluster unable to unmarshal envs - %v", err)
		}

		isClusterFound := false
		// environmentIndex := 0 use this variable
		for _, environment := range envs {
			if environment.ID == environmentID {
				for _, cluster := range environment.Clusters {
					if cluster.ID == clusterID {
						isClusterFound = true
						// TODO REMOVE CLUSTER HERE
						data, err := json.Marshal(envs)
						if err != nil {
							return fmt.Errorf("error deleting cluster unable to marshal envs - %v", err)
						}
						project.Environments = string(data)
						return m.updateProject(accountID, &project)
					}
				}
			}
		}
		if !isClusterFound {
			return fmt.Errorf("error deleting cluster specified cluster not found")
		}
	}
	return fmt.Errorf("error adding cluster project length not equal to one")
}

func (m *Manager) GetCluster(clusterID string) (*model.TableClusters, error) {
	clusters := map[string][]model.TableClusters{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableClusters, utils.OpRead),
		Response: &clusters,
		Params: CrudRequestBody{
			Op:   utils.OpOne,
			Find: &model.TableClusters{ClusterID: clusterID},
		},
	}
	err := utils.HttpRequest(h)
	if err != nil {
		return nil, fmt.Errorf("error getting cluster - %v", err)
	}
	if len(clusters) != 0 {
		return &clusters["result"][0], nil
	}
	return nil, fmt.Errorf("error getting cluster specified cluster not found")
}
