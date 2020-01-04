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

// TODO SHOULD I CHECK IF PROJECT EXISTS OR GRACEFULLY HANDLE IT WHILE PROJECT CREATION OR UPDATION
// UpsertService upserts in services table & calls web hook if upsertion operation is successful
func (m *Manager) UpsertService(ctx context.Context, req *model.Service) error {
	// convert cluster to csv string
	clusters := ""
	for _, cluster := range req.Clusters {
		clusters += fmt.Sprintf("%s,", cluster)
	}
	clusters = strings.TrimSuffix(clusters, ",")

	// store service config as json string in database
	config, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error upserting service unable to marshal - %v", err)
	}

	// send http request to space cloud for db operation
	resp := map[string]interface{}{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableServices, utils.OpUpdate),
		Response: &resp,
		Params: &CrudRequestBody{
			Op: utils.OpUpsert,
			Update: map[string]interface{}{
				"$set": &model.TableServices{
					Environment: req.Environment,
					ProjectID:   req.ProjectID,
					ServiceID:   req.ID,
					Clusters:    clusters,
					Config:      string(config),
					Version:     req.Version,
				},
			},
			Find: &model.TableServices{ProjectID: req.ProjectID, ServiceID: req.ID},
		},
	}
	if err = utils.HttpRequest(h); err != nil {
		return fmt.Errorf("error upserting service - %v %v", err, resp["error"])
	}
	return nil
}

// DeleteService deletes specified service from services table & calls web hook if deletion is successful
func (m *Manager) DeleteService(ctx context.Context, projectID, environmentID, serviceID string) error {
	// send http request to space cloud for db operation
	h := &utils.HttpModel{
		Method: http.MethodPost,
		Url:    getCrudEndpoint(utils.TableServices, utils.OpDelete),
		Params: &CrudRequestBody{
			Op:   utils.OpOne,
			Find: &model.TableServices{ProjectID: projectID, Environment: environmentID, ServiceID: serviceID},
		},
	}
	if err := utils.HttpRequest(h); err != nil {
		return fmt.Errorf("error deleting service - %v", err)
	}
	return nil
}

// TODO for updating service i need to delete the existing service how to do that
// ApplyServiceToClusters applies specified service to all the clusters in service config
func (m *Manager) ApplyServiceToClusters(ctx context.Context, req *model.Service) error {
	for _, clusterID := range req.Clusters {
		// get specified cluster info
		cluster, err := m.GetCluster(clusterID)
		if err != nil {
			return fmt.Errorf("error applying service - %v", err)
		}

		// send http request to space cloud for db operation
		resp := map[string]interface{}{}
		h := &utils.HttpModel{
			Method:   http.MethodPost,
			Url:      fmt.Sprintf("%s/v1/launchpad/service", cluster.Url),
			Params:   req,
			Response: &resp,
		}
		if err = utils.HttpRequest(h); err != nil {
			return fmt.Errorf("error unable to apply service - %v %v", err, resp["error"])
		}
	}

	return nil
}

func (m *Manager) DeleteServiceFromClusters(ctx context.Context, req *model.Service) error {
	for _, clusterID := range req.Clusters {
		// get specified cluster info
		cluster, err := m.GetCluster(clusterID)
		if err != nil {
			return fmt.Errorf("error deleting service from cluster- %v", err)
		}

		// send http request to space cloud for db operation
		resp := map[string]interface{}{}
		h := &utils.HttpModel{
			Method:   http.MethodDelete,
			Url:      fmt.Sprintf("%s/v1/launchpad/service", cluster.Url),
			Params:   req,
			Response: &resp,
		}
		if err = utils.HttpRequest(h); err != nil {
			return fmt.Errorf("error unable to delete service from cluster- %v %v", err, resp["error"])
		}
	}

	return nil
}
