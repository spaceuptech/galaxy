package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleProjectCreation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// defer cancel()
		//
		// req := new(model.ProjectCreateRequest)
		//



		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})

	}
}

func CreateProject() {

}