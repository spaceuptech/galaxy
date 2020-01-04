package file

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

func getCrudEndpoint(tableName, OperationType string) string {
	return fmt.Sprintf("%s/%s/%s", utils.CrudEndpoint, tableName, OperationType)
}

func (m *Manager) addProject(accountID string, req *model.CreateProject) error {
	envs, err := json.Marshal(req.Environments)
	if err != nil {
		return fmt.Errorf("error adding project unable to marshal environments - %v", err)
	}
	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    getCrudEndpoint(utils.TableProjects, utils.OpCreate),
		Params: &CrudRequestBody{
			Op: utils.OpOne,
			Doc: &model.TableProjects{
				ProjectID:    req.ID,
				AccountID:    accountID,
				DefaultEnv:   req.DefaultEnvironment,
				Environments: string(envs),
			}},
	}
	if err = utils.HttpRequest(h); err != nil {
		return fmt.Errorf("error adding project - %v", err)
	}
	return nil
}

func (m *Manager) updateProject(accountID string, req *model.TableProjects) error {
	// TODO CHECK METHOD FOR UPDATE
	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    getCrudEndpoint(utils.TableProjects, utils.OpUpdate),
		Params: &CrudRequestBody{
			Op: utils.OpAll,
			Update: map[string]interface{}{
				"$set": &model.TableProjects{
					Environments: req.Environments,
				},
			},
			Find: &model.TableProjects{
				ProjectID: req.ProjectID,
				AccountID: accountID,
			},
		},
	}
	if err := utils.HttpRequest(h); err != nil {
		return fmt.Errorf("error updating project - %v", err)
	}
	return nil
}
