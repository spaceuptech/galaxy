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

// TODO COMMENT REMAINING
// TODO APPLY LOGRUS COMMANDS EVERTYWHERE
func HandleApplyUIService(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := new(model.Service)
		json.NewDecoder(r.Body).Decode(req)

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		if err := galaxyConfig.UpsertService(ctx, req); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func HandleApplyCLiService(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification TODO WILL THERE BE TOKEN HERE
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		req := new(model.FileStoreEventMessage)
		json.NewDecoder(r.Body).Decode(req)

		if req.Data.Meta.IsDeploy {
			//  TODO PATH VERIFCATION
			if err := galaxyConfig.UpsertService(ctx, req.Data.Meta.Service); err != nil {
				utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func HandleDeleteService(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		vars := mux.Vars(r)

		if err := galaxyConfig.DeleteService(ctx, vars["projectID"], vars["environmentID"], vars["serviceID"]); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

// handles for handling events
func HandleClusterApplyService(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		req := new(model.DatabaseEventMessage)
		json.NewDecoder(r.Body).Decode(req)

		if err := galaxyConfig.ApplyServiceToCluster(ctx, req.Data.Doc); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}

func HandleClusterDeleteService(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		// token verification
		req := new(model.DatabaseEventMessage)
		json.NewDecoder(r.Body).Decode(req)

		if err := galaxyConfig.ApplyServiceToCluster(ctx, req.Data.Doc); err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
		}

		utils.SendEmptySuccessResponse(w, r)
	}
}
