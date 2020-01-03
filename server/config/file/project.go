package file

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

// AddProject adds a project in database if it doesn't exist
func (m *Manager) AddProject(ctx context.Context, req *model.CreateProject) error {
	// check if specified project exists in d
	// atabase
	projects, err := m.GetProject(ctx, req.ID)
	if err != nil {
		return err
	}
	if len(projects) != 0 {
		return fmt.Errorf("error project already exists in database")
	}

	// add project in database
	return m.addProject(req)
}

// GetProject returns the specified project config if it exists in database
func (m *Manager) GetProject(ctx context.Context, projectID string) ([]*model.TableProjects, error) {
	// get specified project from database
	h := &utils.HttpModel{
		Method:  http.MethodGet,
		Url:     getCrudEndpoint(utils.TableProjects, utils.OpRead),
		Headers: nil,
		Params: &CrudRequestBody{
			Op:   utils.OpAll,
			Find: &model.TableProjects{ProjectID: projectID, AccountID: m.accountID}},
		FunctionCallType: utils.SimpleRequest,
	}
	projects, err := utils.HttpRequest(h)
	if err != nil {
		return nil, fmt.Errorf("error getting specified project - %v", err)
	}
	// TODO IN CREATEPROJECT ENV WILL BE STRING FOR DB OPERATION BUT []ENVS FOR REST APIS
	// TODO SHOULD I KEEP THE MARSHALLING RESPONSIBLITY TO FRONTEND ? OR SHOULD I HANDLE IT
	return projects["result"].([]*model.TableProjects), nil
}

// GetProjects return all the projects present in config file
func (m *Manager) GetProjects(ctx context.Context) ([]*model.TableProjects, error) {
	// TODO SHOULD I INTEGRATE GET PROJECT & GET PROJECTS THE DIFFIRENCE IS OF ONLY 1 FIND PARAMTER
	// get all projects of specified account from database
	h := &utils.HttpModel{
		Method:  http.MethodGet,
		Url:     getCrudEndpoint(utils.TableProjects, utils.OpRead),
		Headers: nil,
		Params: &CrudRequestBody{
			Op: utils.OpAll,
			Find: map[string]interface{}{
				utils.ProjectAccount: m.accountID,
			}},
		FunctionCallType: utils.SimpleRequest,
	}
	projects, err := utils.HttpRequest(h)
	if err != nil {
		return nil, fmt.Errorf("error getting projects of specified account - %v", err)
	}

	return projects["result"].([]*model.TableProjects), nil
}

// DeleteProject deletes a specified project from database
func (m *Manager) DeleteProject(ctx context.Context, projectID string) error {
	// TODO WILL THERE BE RACE CONDITION
	// TODO WHAT IF PROJECT DOESN'T EXIST DO I FIRST CHECK THE EXISTENCE THE DELETE ?
	// delete specified project from database
	h := &utils.HttpModel{
		Method:  http.MethodDelete,
		Url:     getCrudEndpoint(utils.TableProjects, utils.OpDelete),
		Headers: nil,
		Params: &CrudRequestBody{
			Op: utils.OpOne,
			Find: map[string]interface{}{
				utils.ProjectID:      projectID,
				utils.ProjectAccount: m.accountID,
			}},
		FunctionCallType: utils.SimpleRequest,
	}
	if _, err := utils.HttpRequest(h); err != nil {
		return fmt.Errorf("error deleting specified project - %v", err)
	}

	return nil
}
