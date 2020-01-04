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

func HandleAddEnvironment(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := new(model.Environment)
		json.NewDecoder(r.Body).Decode(req)

		token, err := auth.VerifyToken(utils.GetTokenFromHeader(r))
		// token verification
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		projectID := mux.Vars(r)["serviceID"]
		if err := galaxyConfig.AddEnvironment(ctx, token["account"].(string), projectID, req); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func HandleDeleteEnvironment(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification
		token, err := auth.VerifyToken(utils.GetTokenFromHeader(r))
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		vars := mux.Vars(r)
		projectID := vars["serviceID"]
		EnvironmentID := vars["environmentID"]

		if err := galaxyConfig.DeleteEnvironment(ctx, token["account"].(string), projectID, EnvironmentID); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func HandleSetDefaultEnvironment(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification
		token, err := auth.VerifyToken(utils.GetTokenFromHeader(r))
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		req := new(model.CreateProject)
		json.NewDecoder(r.Body).Decode(req)

		projectID := mux.Vars(r)["serviceID"]
		if err := galaxyConfig.SetDefaultEnvironment(ctx, token["account"].(string), projectID, req.DefaultEnvironment); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}
