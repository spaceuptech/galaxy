package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils/auth"
)

func HandleCliLogin(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		v := new(model.CliLoginRequest)
		json.NewDecoder(r.Body).Decode(v)

		if !auth.VerifyCliLogin(v.Username, v.Pass) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "unknown username or pass"})
			return
		}

		token, err := auth.GenerateLoginToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token})
		return
	}
}
