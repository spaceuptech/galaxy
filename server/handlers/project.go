package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/server/config"
	"github.com/spaceuptech/launchpad/utils"
	"github.com/spaceuptech/launchpad/utils/auth"
)

func HandleAddProject(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := new(model.CreateProject)
		json.NewDecoder(r.Body).Decode(req)

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		if err := galaxyConfig.AddProject(ctx, req); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func HandleGetProject(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification
		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		projectID := mux.Vars(r)["projectID"]
		project, err := galaxyConfig.GetProject(ctx, projectID)
		if err != nil {
			// TODO WILL THIS CAUSE RACE CONDITION
			utils.SendErrorResponse(w, r, http.StatusNotFound, err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"project": project})
	}
}

func HandleGetProjects(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO TEST VEFIY TOKEN AS PUBLIC KEY WAS NOT GIVEN BUT ERROR WAS NOT THRON PANIC
		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		projects, err := galaxyConfig.GetProjects(ctx)
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

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		projectID := mux.Vars(r)["projectID"]
		if err := galaxyConfig.DeleteProject(ctx, projectID); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}
