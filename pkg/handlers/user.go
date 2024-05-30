package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/models"

	"github.com/go-sql-driver/mysql"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only allow HTTP POST Method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		error_response := map[string]string{"error": fmt.Sprintf("Method not allowed: %s", r.Method)}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	var newUser models.User
	decoder := json.NewDecoder(r.Body)
	// TODO: Add field validation
	if err := decoder.Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error_response := map[string]string{"error": fmt.Sprintf("Error decoding request body: %v", err)}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	result := core.GetDb().Create(&newUser)
	if err := result.Error; err != nil {
		// Check for GORM's duplicate key error
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				// Handle duplicate entry
				w.WriteHeader(http.StatusConflict)
				error_response := map[string]string{"error": fmt.Sprintf("User with this email address already exists: %v", newUser.Email)}
				json.NewEncoder(w).Encode(error_response)
				return
			}
		} else {
			// Handle other potential errors
			w.WriteHeader(http.StatusBadRequest)
			error_response := map[string]string{"error": fmt.Sprintf("Error creating user: %v", err)}
			json.NewEncoder(w).Encode(error_response)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_response := map[string]string{"error": fmt.Sprintf("Error serializing user to JSON: %v", err)}
		json.NewEncoder(w).Encode(error_response)
		return
	}
}
