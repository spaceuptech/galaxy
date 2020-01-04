package file

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

// AddProject adds a project in database if it doesn't exist
func (m *Manager) AddProject(ctx context.Context, accountID string, req *model.CreateProject) error {
	// check if specified project exists in database
	projects, err := m.GetProject(ctx, req.ID, accountID)
	if err != nil {
		return err
	}
	if len(projects) != 0 {
		return fmt.Errorf("error project already exists in database")
	}

	// add project in database
	return m.addProject(accountID, req)
}

// GetProject returns the specified project config if it exists in database
func (m *Manager) GetProject(ctx context.Context, projectID, accountID string) ([]model.TableProjects, error) {
	// get specified project from database
	projects := map[string][]model.TableProjects{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableProjects, utils.OpRead),
		Response: &projects,
		Params: &CrudRequestBody{
			Op:   utils.OpAll,
			Find: &model.TableProjects{ProjectID: projectID, AccountID: accountID}},
	}
	err := utils.HttpRequest(h)
	if err != nil {
		return nil, fmt.Errorf("error getting specified project - %v %v", err, projects["error"])
	}
	// TODO IN CREATEPROJECT ENV WILL BE STRING FOR DB OPERATION BUT []ENVS FOR REST APIS
	// TODO SHOULD I KEEP THE MARSHALLING RESPONSIBLITY TO FRONTEND ? OR SHOULD I HANDLE IT
	return projects["result"], nil
}

// GetProjects return all the projects present in config file
func (m *Manager) GetProjects(ctx context.Context, accountID string) ([]model.TableProjects, error) {
	// TODO SHOULD I INTEGRATE GET PROJECT & GET PROJECTS THE DIFFIRENCE IS OF ONLY 1 FIND PARAMTER
	// get all projects of specified account from database
	projects := map[string][]model.TableProjects{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Response: &projects,
		Url:      getCrudEndpoint(utils.TableProjects, utils.OpRead),
		Params: &CrudRequestBody{
			Op: utils.OpAll,
			Find: map[string]interface{}{
				utils.ProjectAccount: accountID,
			}},
	}
	err := utils.HttpRequest(h)
	if err != nil {
		return nil, fmt.Errorf("error getting projects of specified account - %v -%v", err, projects["error"])
	}

	return projects["result"], nil
}

// DeleteProject deletes a specified project from database
func (m *Manager) DeleteProject(ctx context.Context, projectID, accountID string) error {
	// TODO WILL THERE BE RACE CONDITION
	// TODO WHAT IF PROJECT DOESN'T EXIST DO I FIRST CHECK THE EXISTENCE THE DELETE ?
	// delete specified project from database
	h := &utils.HttpModel{
		Method:  http.MethodPost,
		Url:     getCrudEndpoint(utils.TableProjects, utils.OpDelete),
		Headers: nil,
		Params: &CrudRequestBody{
			Op: utils.OpOne,
			Find: map[string]interface{}{
				utils.ProjectID:      projectID,
				utils.ProjectAccount: accountID,
			}},
	}
	if err := utils.HttpRequest(h); err != nil {
		return fmt.Errorf("error deleting specified project - %v", err)
	}

	return nil
}

// ApplyServiceToClusters applies specified service to all the clusters in service config
func (m *Manager) CreateProjectInClusters(ctx context.Context, req *model.CreateProject) error {
	// for _, clusterID := range req.Clusters {
	// 	// get specified cluster info
	// 	cluster, err := m.GetCluster(clusterID)
	// 	if err != nil {
	// 		return fmt.Errorf("error applying service - %v", err)
	// 	}
	//
	// 	// send http request to space cloud for db operation
	// 	resp := map[string]interface{}{}
	// 	h := &utils.HttpModel{
	// 		Method:   http.MethodPost,
	// 		Url:      fmt.Sprintf("%s/v1/launchpad/service", cluster.Url),
	// 		Params:   req,
	// 		Response: &resp,
	// 	}
	// 	if err = utils.HttpRequest(h); err != nil {
	// 		return fmt.Errorf("error unable to apply service - %v %v", err, resp["error"])
	// 	}
	// }

	return nil
}

func (m *Manager) DeleteProjectFromClusters(ctx context.Context, req *model.CreateProject) error {
	// for _, clusterID := range req.Clusters {
	// 	// get specified cluster info
	// 	cluster, err := m.GetCluster(clusterID)
	// 	if err != nil {
	// 		return fmt.Errorf("error deleting service from cluster- %v", err)
	// 	}
	//
	// 	// send http request to space cloud for db operation
	// 	resp := map[string]interface{}{}
	// 	h := &utils.HttpModel{
	// 		Method:   http.MethodDelete,
	// 		Url:      fmt.Sprintf("%s/v1/launchpad/service", cluster.Url),
	// 		Params:   req,
	// 		Response: &resp,
	// 	}
	// 	if err = utils.HttpRequest(h); err != nil {
	// 		return fmt.Errorf("error unable to delete service from cluster- %v %v", err, resp["error"])
	// 	}
	// }
	//
	return nil
}
