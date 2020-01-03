package file

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

// TODO SET HEADERS OF HTTP REQUEST TO APPLICATION JSON ?
func (m *Manager) UpsertService(ctx context.Context, req *model.Service) error {
	// convert cluster to csv string
	clusters := "" // csv of clusters
	for _, cluster := range req.Clusters {
		clusters += fmt.Sprintf("%s,", cluster)
	}
	clusters = strings.TrimSuffix(clusters, ",")

	config, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error upserting service unable to marshal - %v", err)
	}
	// TODO IS CREATE ENDPOINT WILL UPSERT ?
	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    getCrudEndpoint(utils.TableServices, utils.OpCreate),
		Params: &CrudRequestBody{
			Op: utils.OpUpsert,
			Doc: &model.TableService{
				Environment: req.Environment,
				ProjectID:   req.ProjectID,
				ServiceID:   req.ID,
				Clusters:    clusters,
				Config:      string(config),
				Version:     req.Version,
			},
			Find: map[string]interface{}{"service_id": req.ID},
		},
		FunctionCallType: utils.SimpleRequest,
	}
	_, err = utils.HttpRequest(h)
	if err != nil {
		return fmt.Errorf("error upserting service - %v", err)
	}
	return nil
}

func (m *Manager) DeleteService(ctx context.Context, serviceID string) error {
	h := &utils.HttpModel{
		Method: http.MethodDelete,
		Url:    getCrudEndpoint(utils.TableServices, utils.OpDelete),
		Params: &CrudRequestBody{
			Op:   utils.OpOne,
			Find: &model.TableService{ServiceID: serviceID},
		},
		FunctionCallType: utils.SimpleRequest,
	}
	_, err := utils.HttpRequest(h)
	if err != nil {
		return fmt.Errorf("error deleting service - %v", err)
	}
	return nil
}

func (m *Manager) ApplyService(ctx context.Context, req *model.Service) error {
	cluster, err := m.GetCluster(req.ID)
	if err != nil {
		return fmt.Errorf("error applying service - %v", err)
	}

	h := &utils.HttpModel{
		Method:           http.MethodPost,
		Url:              cluster.Url,
		Params:           req,
		FunctionCallType: utils.SimpleRequest,
	}
	_, err = utils.HttpRequest(h)
	if err != nil {
		fmt.Errorf("error unable to apply service - %v", err)
	}
	return nil
}

// TODO I GUESS THIS FUNCTION IS NOT USED
// GetProject returns the specified project config if it exists in database
func (m *Manager) getService(ctx context.Context, projectID string) ([]*model.TableProjects, error) {
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
