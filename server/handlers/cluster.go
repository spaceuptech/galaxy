package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/server/config"
	"github.com/spaceuptech/launchpad/utils"
	"github.com/spaceuptech/launchpad/utils/auth"
)

func HandleClusterRegistration(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// if _, err := auth.VerifyToken(utils.GetToken(r)); err != nil {
		// 	utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		// }

		request := new(model.RegisterClusterRequest)
		json.NewDecoder(r.Body).Decode(request)

		// TODO WHAT IF THE CLUSTER ALREADY EXIST IN DATABASE
		log.Println("updating custer")
		if _, err := updateCluster(ctx, request, utils.ClusterAlive); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}
		log.Println("updated custer")

		// TODO PING THE REGISTERED CLUSTER TO UPDATE THE STATUS IN DATABASE
		ticker := time.NewTicker(3 * time.Hour)
		done := make(chan bool)

		go func() {
			clusterAliveCount := 1
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					h := &utils.HttpModel{
						Method:           http.MethodPost,
						Url:              request.Url,
						FunctionCallType: utils.SimpleRequest,
					}
					_, err := utils.HttpRequest(h)
					if err != nil {
						if clusterAliveCount == utils.MaximumPingRetries {
							// TODO UPDATE THE CLUSTER STATUS TO DEAD IN DATABASE
							if _, err := updateCluster(ctx, request, utils.ClusterDead); err != nil {
								utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
							}
							// TODO WHAT IF UPDATION FAILED
							ticker.Stop()
							done <- true
						}
						clusterAliveCount++
					}
				}
			}
		}()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})

	}
}

func updateCluster(ctx context.Context, request *model.RegisterClusterRequest, status string) (interface{}, error) {
	return utils.FireGraphqlQuery(&model.InsertRequest{
		Query: utils.UpsertInClusterTable,
		Variables: map[string]interface{}{
			"cluster_id":  request.ClusterID,
			"runner_type": request.RunnerType,
			"status":      status,
			"url":         request.Url,
		},
	}, utils.GraphqlMutation)
}

func HandleAddProjectCluster(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := new(model.Cluster)
		json.NewDecoder(r.Body).Decode(req)

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		vars := mux.Vars(r)
		projectID := vars["projectID"]
		EnvironmentID := vars["environmentID"]

		if err := galaxyConfig.AddCluster(ctx, projectID, EnvironmentID, req); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func HandleDeleteProjectCluster(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		vars := mux.Vars(r)
		projectID := vars["projectID"]
		EnvironmentID := vars["environmentID"]
		clusterID := vars["clusterID"]

		if err := galaxyConfig.DeleteCluster(ctx, projectID, EnvironmentID, clusterID); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}
