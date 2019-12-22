package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.WriteHeader(http.StatusUpgradeRequired)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "upgrade is required for login"})
	}
}
