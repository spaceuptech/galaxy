package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/server/config"
	"github.com/spaceuptech/launchpad/utils"
	"github.com/spaceuptech/launchpad/utils/auth"
)

// TODO HMAC TOKEN
func HandleLogin(auth *auth.Module, galaxyConfig *config.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		v := new(model.LoginPayload)
		json.NewDecoder(r.Body).Decode(v)

		if !auth.VerifyCliLogin(v.Username, v.Key) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("unknown username or key")})
			return
		}
		// TODO SEND PROJECT IN RESPONSE
		token, err := auth.GenerateLoginToken()
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		projects, err := galaxyConfig.GetProjects(ctx, v.Username)
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		scToken, err := auth.GenerateHS256Token()
		if err != nil {
			utils.SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "projects": projects, "fileToken": scToken})
		return
	}
}
