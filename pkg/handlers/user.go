package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Only allow HTTP POST Method
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errorResponse := map[string]string{"error": fmt.Sprintf("Method not allowed: %s", r.Method)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var newUser models.User
	decoder := json.NewDecoder(r.Body)
	// TODO: Add field validation
	if err := decoder.Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := map[string]string{"error": fmt.Sprintf("Error decoding request body: %v", err)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Hash the password before saving to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": fmt.Sprintf("Error hashing password: %v", err)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	newUser.Password = string(hashedPassword)

	result := core.GetDb().Create(&newUser)
	if err := result.Error; err != nil {
		// Check for GORM's duplicate key error
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				// Handle duplicate entry
				w.WriteHeader(http.StatusConflict)
				errorResponse := map[string]string{"error": fmt.Sprintf("User with this email address already exists: %v", newUser.Email)}
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
		} else {
			// Handle other potential errors
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := map[string]string{"error": fmt.Sprintf("Error creating user: %v", err)}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": fmt.Sprintf("Error serializing user to JSON: %v", err)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	// Only allow HTTP POST Method
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errorResponse := map[string]string{"error": fmt.Sprintf("Method not allowed: %s", r.Method)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var loginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginPayload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := map[string]string{"error": fmt.Sprintf("Error decoding request body: %v", err)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var user models.User
	result := core.GetDb().Where("email = ?", loginPayload.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			error_response := map[string]string{"error": "Invalid credentials"}
			json.NewEncoder(w).Encode(error_response)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			error_response := map[string]string{"error": fmt.Sprintf("Error retrieving user: %v", result.Error)}
			json.NewEncoder(w).Encode(error_response)
			return
		}
	}

	// Compare the hashed password stored in the database with the hash of the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		error_response := map[string]string{"error": "Invalid credentials"}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	// Create a response object with a "token" key
	encoder := json.NewEncoder(w)
	response := struct {
		Token string `json:"token"`
	}{
		Token: "test_token",
	}

	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := map[string]string{"error": fmt.Sprintf("Error serializing user to JSON: %v", err)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
}
