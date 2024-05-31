package handlers

import (
	"encoding/json"
	"fmt"
	"math"
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
		w.WriteHeader(http.StatusNotFound)
		error_response := map[string]string{"error": "User not found in context"}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	// Only allow HTTP POST Method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errorResponse := map[string]string{"error": fmt.Sprintf("Method not allowed: %s", r.Method)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var swipePayload struct {
		TargetID  uint64 `json:"targetID"`
		SwipeType string `json:"swipeType"`
	}

	err := json.NewDecoder(r.Body).Decode(&swipePayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error_response := map[string]string{"error": "Payload is not valid"}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	// Check if the tragetID is the same as swiperID
	if swipePayload.TargetID == contextUser.ID {
		w.WriteHeader(http.StatusBadRequest)
		error_response := map[string]string{"error": "Cannot swipe on yourself"}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	// Check if the tragetID user exists
	var targetUser models.User
	err = core.GetDb().First(&targetUser, swipePayload.TargetID).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error_response := map[string]string{"error": "Target user not found"}
		json.NewEncoder(w).Encode(error_response)
		return
	}

	swipe := models.Swipe{
		SwiperID:  contextUser.ID,
		TargetID:  swipePayload.TargetID,
		SwipeType: swipePayload.SwipeType,
	}

	core.GetDb().Create(&swipe)
	if swipe.SwipeType == "YES" {

		// Increase the Target User Total Likes
		targetUser.TotalLikesReceived++
		targetUser.AttractivenessScore = calculateAttractivenessScore(targetUser.TotalLikesReceived, targetUser.TotalDislikesReceived)

		// Save the changes to the target user
		core.GetDb().Save(&targetUser)

		var targetSwipe models.Swipe
		err := core.GetDb().Where("swiper_id = ? AND target_id = ? AND swipe_type = ?", swipe.TargetID, swipe.SwiperID, "YES").First(&targetSwipe).Error

		// If both users swipped YES on each other
		if err == nil {
			match := models.Match{
				User1ID: swipe.SwiperID,
				User2ID: swipe.TargetID,
			}
			core.GetDb().Create(&match)
			encoder := json.NewEncoder(w)
			createdUserResponse := UserSwipeMatchedResonse{
				Matched: true,
				MatchID: match.ID,
			}

			if err := encoder.Encode(createdUserResponse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorResponse := map[string]string{"error": fmt.Sprintf("Error serializing UserSwipeMatchedResonse to JSON: %v", err)}
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
			return
		}
	} else {
		// Increase the Target User Total Dislikes
		targetUser.TotalDislikesReceived++
		targetUser.AttractivenessScore = calculateAttractivenessScore(targetUser.TotalLikesReceived, targetUser.TotalDislikesReceived)
		// Save the changes to the target user
		core.GetDb().Save(&targetUser)
	}

	encoder := json.NewEncoder(w)
	createdUserResponse := UserSwipeNotMatchedResonse{
		Matched: false,
	}

	if err := encoder.Encode(createdUserResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": fmt.Sprintf("Error serializing UserSwipeNotMatchedResonse to JSON: %v", err)}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
}

func calculateAttractivenessScore(likes, dislikes int) float64 {
	totalSwipes := likes + dislikes
	if totalSwipes == 0 {
		return 0.0
	}
	score := float64(likes) / float64(totalSwipes)

	// Round to the nearest 2 digits
	return math.Round(score*100) / 100
}
