package file

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

// AddProject is used to add specified project in projects table if it doesn't exist
func (m *Manager) AddProject(ctx context.Context, accountID string, req *model.CreateProject) error {
	// get specified project info from projects table
	projects, err := m.GetProject(ctx, req.ID, accountID)
	if err != nil {
		return err
	}

	if len(projects) != 0 {
		logrus.Errorf("error project already exists in database")
		return fmt.Errorf("error project already exists in database")
	}

	// add project in projects table
	return m.addProject(accountID, req)
}

// GetProject is used to get specified project config if it exists in database
func (m *Manager) GetProject(ctx context.Context, projectID, accountID string) ([]model.TableProjects, error) {
	// send http request to space cloud for db operation
	projects := map[string][]model.TableProjects{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableProjects, utils.OpRead),
		Response: &projects,
		Params: &model.CrudRequestPayload{
			Op:   utils.OpAll,
			Find: &model.TableProjects{ProjectID: projectID, AccountID: accountID}},
	}
	err := utils.HttpRequest(h)
	if err != nil {
		logrus.Errorf("error getting specified project - %v %v", err, projects["error"])
		return nil, err
	}
	return projects["result"], nil
}

// GetProjects is used to get all the projects in projects table for specified  account
func (m *Manager) GetProjects(ctx context.Context, accountID string) ([]model.TableProjects, error) {
	// send http request to space cloud for db operation
	projects := map[string][]model.TableProjects{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Response: &projects,
		Url:      getCrudEndpoint(utils.TableProjects, utils.OpRead),
		Params: &model.CrudRequestPayload{
			Op: utils.OpAll,
			Find: map[string]interface{}{
				utils.ProjectAccount: accountID,
			}},
	}
	err := utils.HttpRequest(h)
	if err != nil {
		logrus.Errorf("error getting projects of specified account - %v -%v", err, projects["error"])
		return nil, err
	}
	return projects["result"], nil
}

// DeleteProject deletes a specified project from database
func (m *Manager) DeleteProject(ctx context.Context, projectID, accountID string) error {
	// TODO WILL THERE BE RACE CONDITION
	// send http request to space cloud for db operation
	resp := map[string]interface{}{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableProjects, utils.OpDelete),
		Response: &resp,
		Params: &model.CrudRequestPayload{
			Op: utils.OpOne,
			Find: map[string]interface{}{
				utils.ProjectID:      projectID,
				utils.ProjectAccount: accountID,
			}},
	}
	if err := utils.HttpRequest(h); err != nil {
		logrus.Errorf("error deleting specified project - %v", err)
		return err
	}

	return nil
}

// ApplyServiceToClusters is used to apply specified service to all the clusters in service config
func (m *Manager) CreateProjectInClusters(ctx context.Context, req *model.CreateProject) error {
	for clusterID, environments := range createStructure(req) {
		// get specified cluster info from clusters table
		cluster, err := m.GetCluster(clusterID)
		if err != nil {
			logrus.Errorf("error creating project on cluster - %v", err)
			return err
		}

		// send http request to space cloud for db operation
		resp := map[string]interface{}{}
		h := &utils.HttpModel{
			Method: http.MethodPost,
			Url:    fmt.Sprintf("%s/v1/launchpad/project", cluster.Url),
			Params: &model.CreateClusterProjectPayload{
				ProjectID:    req.ID,
				Environments: environments,
			},
			Response: &resp,
		}
		if err = utils.HttpRequest(h); err != nil {
			logrus.Errorf("error unable to create project on cluster - %v %v", err, resp["error"])
			return err
		}
	}
	return nil
}

// DeleteProjectFromClusters is used to delete specified cluster from all the clusters present in service config
func (m *Manager) DeleteProjectFromClusters(ctx context.Context, req *model.CreateProject) error {
	for clusterID, environments := range createStructure(req) {
		// get specified cluster info
		cluster, err := m.GetCluster(clusterID)
		if err != nil {
			logrus.Errorf("error creating project on cluster - %v", err)
			return err
		}

		// send http request to space cloud for db operation
		resp := map[string]interface{}{}
		h := &utils.HttpModel{
			Method: http.MethodDelete,
			Url:    fmt.Sprintf("%s/v1/launchpad/project", cluster.Url),
			Params: &model.CreateClusterProjectPayload{
				ProjectID:    req.ID,
				Environments: environments,
			},
			Response: &resp,
		}
		if err = utils.HttpRequest(h); err != nil {
			logrus.Errorf("error unable to create project on cluster - %v %v", err, resp["error"])
			return err
		}
	}
	return nil
}

// createStructure generates a data structure which helps creating/deleting service info in/from all the cluster present in service config
func createStructure(req *model.CreateProject) map[string][]model.Environment {
	clusters := map[string][]model.Environment{}
	for _, environment := range req.Environments {
		for _, cluster := range environment.Clusters {
			arrValue, ok := clusters[cluster.ID]
			if ok {
				clusters[cluster.ID] = append(arrValue, *environment)
			}
			clusters[cluster.ID] = []model.Environment{*environment}
		}
	}
	return clusters
}
