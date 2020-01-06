package file

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/opencontainers/runc/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

// AddProjectCluster is used to add specified cluster to the projects config it doesn't exits
func (m *Manager) AddProjectCluster(ctx context.Context, accountID, projectID, environmentID string, req *model.Cluster) error {
	// get specified project info from projects table
	projects, err := m.GetProject(ctx, projectID, accountID)
	if err != nil {
		logrus.Errorf("error adding cluster - %v", err)
		return err
	}

	// length should be 1 as project id stored in table is unique
	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		// convert to envs as it is stored as json string in projects table
		if err := json.Unmarshal([]byte(project.Environments), &envs); err != nil {
			logrus.Errorf("error adding cluster unable to unmarshal envs - %v", err)
			return err
		}

		// check if cluster already exists
		isClusterFound, environmentIndex := false, 0
		for _, environment := range envs {
			if environment.ID == environmentID {
				for _, cluster := range environment.Clusters {
					isClusterFound = true
					if cluster.ID == req.ID {
						logrus.Errorf("error adding cluster specified cluster already exists")
						return err
					}
				}
			}
		}

		// if doesn't exits then add new cluster & update the table
		if isClusterFound {
			envs[environmentIndex].Clusters = append(envs[environmentIndex].Clusters, req)
			// convert envs to json string as it stored as json string in table
			data, err := json.Marshal(envs)
			if err != nil {
				logrus.Errorf("error adding cluster unable to marshal envs - %v", err)
				return err
			}
			project.Environments = string(data)
			return m.updateProject(accountID, &project)
		}
	}
	logrus.Errorf("error adding cluster project length not equal to one")
	return err
}

// DeleteProjectCluster is used to remove specified cluster from the project config if it exists
func (m *Manager) DeleteProjectCluster(ctx context.Context, accountID, projectID, environmentID, clusterID string) error {
	// get specified project info from projects table
	projects, err := m.GetProject(ctx, projectID, accountID)
	if err != nil {
		logrus.Errorf("error adding cluster - %v", err)
		return err
	}

	// length should be 1 as project id stored in table is unique
	if len(projects) == 1 {
		project := projects[0]
		envs := []*model.Environment{}
		// convert to envs as it is stored as json string in projects table
		if err := json.Unmarshal([]byte(project.Environments), &envs); err != nil {
			logrus.Errorf("error adding cluster unable to unmarshal envs - %v", err)
			return err
		}

		// check if specified cluster exits
		isClusterFound := false
		for environmentIndex, environment := range envs {
			if environment.ID == environmentID {
				for clusterIndex, cluster := range environment.Clusters {
					if cluster.ID == clusterID {
						isClusterFound = true
						envs[environmentIndex].Clusters = removeClusterAtIndex(environment.Clusters, clusterIndex)
						// convert envs to json string before updating table
						data, err := json.Marshal(envs)
						if err != nil {
							logrus.Errorf("error deleting cluster unable to marshal envs - %v", err)
							return err
						}
						project.Environments = string(data)
						return m.updateProject(accountID, &project)
					}
				}
			}
		}

		// throw error if specified cluster doesn't exits
		if !isClusterFound {
			logrus.Errorf("error deleting cluster specified cluster not found")
			return err
		}
	}
	logrus.Errorf("error deleting cluster project length not equal to one")
	return err
}

// GetCluster is used to get specified cluster info from clusters table
func (m *Manager) GetCluster(clusterID string) (*model.TableClusters, error) {
	// send http request to space cloud for db operation
	clusters := map[string][]model.TableClusters{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableClusters, utils.OpRead),
		Response: &clusters,
		Params: model.CrudRequestPayload{
			Op:   utils.OpOne,
			Find: &model.TableClusters{ClusterID: clusterID},
		},
	}
	err := utils.HttpRequest(h)
	if err != nil {
		logrus.Errorf("error getting cluster - %v %v", err, clusters["error"])
		return nil, err
	}
	if len(clusters) != 0 {
		return &clusters["result"][0], nil
	}
	logrus.Errorf("error getting cluster specified cluster not found")
	return nil, err
}

// UpdateCluster is used to update registered clusters in clusters table
func (m *Manager) UpdateCluster(ctx context.Context, request *model.RegisterClusterPayload, status string) error {
	// send http request to space cloud for db operation
	resp := map[string]interface{}{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableClusters, utils.OpUpdate),
		Response: &resp,
		Params: &model.CrudRequestPayload{
			Op:   utils.OpUpsert,
			Find: nil,
			Update: &model.TableClusters{
				ClusterID:  request.ClusterID,
				RunnerType: request.RunnerType,
				Status:     status,
				Url:        request.Url,
			},
		},
	}
	if err := utils.HttpRequest(h); err != nil {
		logrus.Errorf("error updating cluster %v", resp["error"])
		return err
	}
	return nil
}
