package file

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/utils"
)

// getCrudEndpoint is used to get CRUD endpoint of space cloud
func getCrudEndpoint(tableName, OperationType string) string {
	return fmt.Sprintf("%s/%s/%s", utils.CrudEndpoint, tableName, OperationType)
}

// addProject is used to add specified project in projects table
func (m *Manager) addProject(accountID string, req *model.CreateProject) error {
	// convert envs to json string as it stored as json string in projects table
	envs, err := json.Marshal(req.Environments)
	if err != nil {
		logrus.Errorf("error adding project unable to marshal environments - %v", err)
		return err
	}

	resp := map[string]interface{}{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableProjects, utils.OpCreate),
		Response: &resp,
		Params: &model.CrudRequestPayload{
			Op: utils.OpOne,
			Doc: &model.TableProjects{
				ProjectID:    req.ID,
				AccountID:    accountID,
				DefaultEnv:   req.DefaultEnvironment,
				Environments: string(envs),
			}},
	}
	if err = utils.HttpRequest(h); err != nil {
		logrus.Errorf("error adding project - %v %v", err, resp["error"])
		return err
	}
	return nil
}

// updateProject is used to update specified project
func (m *Manager) updateProject(accountID string, req *model.TableProjects) error {
	resp := map[string]interface{}{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableProjects, utils.OpUpdate),
		Response: &resp,
		Params: &model.CrudRequestPayload{
			Op: utils.OpUpsert,
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
		logrus.Errorf("error updating project - %v %v", err, resp["error"])
		return err
	}
	return nil
}

// removeEnvironmentAtIndex removes environment at specified index of project config
func removeEnvironmentAtIndex(s []*model.Environment, i int) []*model.Environment {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// removeClusterAtIndex removes cluster at specified index of project config
func removeClusterAtIndex(s []*model.Cluster, i int) []*model.Cluster {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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
