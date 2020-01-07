package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/spaceuptech/galaxy/utils"
	"github.com/spaceuptech/galaxy/utils/auth"
)

// HandleProvidePublicKey sends public key to runner
func HandleProvidePublicKey(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// token verification
		if _, err := auth.VerifyToken(utils.GetTokenFromHeader(r)); err != nil {
			utils.SendErrorResponse(w, r, http.StatusUnauthorized, err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"pem": auth.GetPublicKey()})
	}
}
