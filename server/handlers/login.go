package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils/auth"
)

func HandleLogin(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		v := new(model.CliLoginRequest)
		json.NewDecoder(r.Body).Decode(v)

		if !auth.VerifyCliLogin(v.Username, v.Key) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": fmt.Sprintf("unknown username or key")})
			return
		}
		// TODO SEND PROJECT IN RESPONSE
		token, err := auth.GenerateLoginToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
			return
		}
		// log.Println("generating token", token)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token})
		return
	}
}
