package file

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/opencontainers/runc/Godeps/_workspace/src/github.com/Sirupsen/logrus"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/utils"
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
		logrus.Errorf("error upserting service unable to marshal - %v", err)
		return err
	}

	// send http request to space cloud for db operation
	resp := map[string]interface{}{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableServices, utils.OpUpdate),
		Response: &resp,
		Params: &model.CrudRequestPayload{
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
			Find: &model.TableServices{ProjectID: req.ProjectID, ServiceID: req.ID, Environment: req.Environment, Version: req.Version},
		},
	}
	if err = utils.HttpRequest(h); err != nil {
		logrus.Errorf("error upserting service - %v error message %v", err, resp["error"])
		return err
	}
	return nil
}

// DeleteService deletes specified service from services table & calls web hook if deletion is successful
func (m *Manager) DeleteService(ctx context.Context, projectID, environmentID, serviceID, version string) error {
	// send http request to space cloud for db operation
	resp := map[string]interface{}{}
	h := &utils.HttpModel{
		Method:   http.MethodPost,
		Url:      getCrudEndpoint(utils.TableServices, utils.OpDelete),
		Response: &resp,
		Params: &model.CrudRequestPayload{
			Op:   utils.OpOne,
			Find: &model.TableServices{ProjectID: projectID, ServiceID: serviceID, Environment: environmentID, Version: version},
		},
	}
	if err := utils.HttpRequest(h); err != nil {
		logrus.Errorf("error deleting service - %v %v", err, resp["error"])
		return err
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
			logrus.Errorf("error applying specified service to cluster - %v", err)
			return err
		}

		// send http request to space cloud for db operation
		resp := map[string]interface{}{}
		h := &utils.HttpModel{
			Method:   http.MethodPost,
			Url:      fmt.Sprintf("%s/v1/galaxy/service", cluster.Url),
			Params:   req,
			Response: &resp,
		}
		if err = utils.HttpRequest(h); err != nil {
			logrus.Errorf("error unable to apply specified service to cluster - %v %v", err, resp["error"])
			return err
		}
	}

	return nil
}

// DeleteServiceFromClusters deletes specified service from clusters present service config
func (m *Manager) DeleteServiceFromClusters(ctx context.Context, req *model.Service) error {
	for _, clusterID := range req.Clusters {
		// get specified cluster info
		cluster, err := m.GetCluster(clusterID)
		if err != nil {
			logrus.Errorf("error deleting service from cluster- %v", err)
			return err
		}

		// send http request to space cloud for db operation
		resp := map[string]interface{}{}
		h := &utils.HttpModel{
			Method:   http.MethodDelete,
			Url:      fmt.Sprintf("%s/v1/galaxy/service", cluster.Url),
			Params:   req,
			Response: &resp,
		}
		if err = utils.HttpRequest(h); err != nil {
			logrus.Errorf("error unable to delete service from cluster- %v %v", err, resp["error"])
			return err
		}
	}

	return nil
}
