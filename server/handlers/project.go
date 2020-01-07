package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/server/config"
	"github.com/spaceuptech/galaxy/utils"
	"github.com/spaceuptech/galaxy/utils/auth"
)

func HandleAddProject(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := new(model.CreateProject)
		json.NewDecoder(r.Body).Decode(req)

		// token verification
		token, err := auth.VerifyToken(utils.GetTokenFromHeader(r))
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		if err := galaxyConfig.AddProject(ctx, token["account"].(string), req); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

// TODO RETURN AFFTER SENDING SUCCES RESPONSE
func HandleGetProject(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification
		// token verification
		token, err := auth.VerifyToken(utils.GetTokenFromHeader(r))
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		projectID := mux.Vars(r)["serviceID"]
		projects, err := galaxyConfig.GetProject(ctx, token["account"].(string), projectID)
		if err != nil {
			// TODO WILL THIS CAUSE RACE CONDITION
			utils.SendErrorResponse(w, r, http.StatusNotFound, err)
			return
		}

		if len(projects) == 0 {
			utils.SendErrorResponse(w, r, http.StatusNotFound, fmt.Errorf("error specified project doesn't exist in database"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"project": projects[0]})
	}
}

func HandleGetProjects(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO TEST VEFIY TOKEN AS PUBLIC KEY WAS NOT GIVEN BUT ERROR WAS NOT THRON PANIC
		// token verification
		token, err := auth.VerifyToken(utils.GetTokenFromHeader(r))
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		projects, err := galaxyConfig.GetProjects(ctx, token["account"].(string))
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"projects": projects})
	}
}

func HandleDeleteProject(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//
		// token verification
		token, err := auth.VerifyToken(utils.GetTokenFromHeader(r))
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		projectID := mux.Vars(r)["serviceID"]
		if err := galaxyConfig.DeleteProject(ctx, token["account"].(string), projectID); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}
