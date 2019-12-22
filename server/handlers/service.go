package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleServiceCreation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// TODO AUTHORIZATION CHECK LIKE TOKEN VALIDATION

		// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// defer cancel()



		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})

	}
}
