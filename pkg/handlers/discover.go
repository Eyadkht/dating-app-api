package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/models"
)

type PotentialMatchesResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

func GetPotentialMatches(w http.ResponseWriter, r *http.Request) {

	// Set content-type to json
	w.Header().Set("Content-Type", "application/json")

	// Retrieve user from context
	contextUser, ok := r.Context().Value(core.UserContextKey).(models.User)
	if !ok {
		// Handle the case where user is not found in context (unexpected)
		w.WriteHeader(http.StatusInternalServerError)
		error_response := map[string]string{"error": "Error: User not found in context"}
		json.NewEncoder(w).Encode(error_response)
		return
	}
	fmt.Println(contextUser)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed: %s", r.Method)
		error_response := map[string]string{"error": fmt.Sprintf("Method not allowed: %s", r.Method)}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	// Fetch all users from the database
	users := []models.User{}
	result := core.GetDb().Omit("password", "email", "Token").Find(&users)

	// Convert User slices to PublicUserResponse slices
	// Used as a data transfer object to omit Token and Password fields
	potentialMatches := make([]PotentialMatchesResponse, len(users))
	for i, user := range users {
		potentialMatches[i] = PotentialMatchesResponse{
			ID:     user.ID,
			Name:   user.Name,
			Gender: user.Gender,
			Age:    user.Age,
		}
	}

	if err := result.Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_response := map[string]string{"error": "Error fetching users"}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	encoder := json.NewEncoder(w)

	// Create a response object with a "results" key
	response := struct {
		Results []PotentialMatchesResponse `json:"results"`
	}{
		Results: potentialMatches,
	}
	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error_response := map[string]string{"error": fmt.Sprintf("Error serializing users to JSON: %v", err)}
		json.NewEncoder(w).Encode(error_response)
		return
	}
}
