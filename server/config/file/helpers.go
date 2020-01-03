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

func (m *Manager) addProject(req *model.CreateProject) error {
	envs, err := json.Marshal(req.Environments)
	if err != nil {
		return fmt.Errorf("error adding project unable to marshal environments - %v", err)
	}
	h := &utils.HttpModel{
		Method:  http.MethodPost,
		Url:     getCrudEndpoint(utils.TableProjects, utils.OpCreate),
		Headers: nil,
		Params: &CrudRequestBody{
			Op: utils.OpOne,
			Doc: map[string]interface{}{
				utils.ProjectID:         req.ID,
				utils.ProjectAccount:    m.accountID,
				utils.ProjectDefaultEnv: req.DefaultEnvironment,
				utils.ProjectEnvs:       envs,
			}},
		FunctionCallType: utils.SimpleRequest,
	}
	_, err = utils.HttpRequest(h)
	if err != nil {
		return fmt.Errorf("error adding project - %v", err)
	}
	return nil
}

func (m *Manager) updateProject(req *model.TableProjects) error {
	// TODO CHECK METHOD FOR UPDATE
	h := &utils.HttpModel{
		Method:  http.MethodPost,
		Url:     getCrudEndpoint(utils.TableProjects, utils.OpUpdate),
		Headers: nil,
		Params: &CrudRequestBody{
			Op:  utils.OpOne,
			Doc: req,
		},
		FunctionCallType: utils.SimpleRequest,
	}
	_, err := utils.HttpRequest(h)
	if err != nil {
		return fmt.Errorf("error updating project - %v", err)
	}
	return nil
}
