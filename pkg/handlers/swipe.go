package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/models"
)

type UserSwipeMatchedResonse struct {
	Matched bool   `json:"matched"`
	MatchID uint64 `json:"matchID"`
}

type UserSwipeNotMatchedResonse struct {
	Matched bool `json:"matched"`
}

func UserSwipe(w http.ResponseWriter, r *http.Request) {

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
	fmt.Println(contextUser.ID)

	// Only allow HTTP POST Method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errorResponse := map[string]string{"error": fmt.Sprintf("Method not allowed: %s", r.Method)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	// Used as a data transfer object to omit the Token field
	createdUserResponse := UserSwipeMatchedResonse{
		Matched: false,
		MatchID: 2,
	}
	if err := encoder.Encode(createdUserResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": fmt.Sprintf("Error serializing user to JSON: %v", err)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
}
