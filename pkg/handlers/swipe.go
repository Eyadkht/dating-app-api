package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/models"
	"muzz-dating/pkg/utils"
)

type UserSwipeMatchedResonse struct {
	Matched bool   `json:"matched"`
	MatchID uint64 `json:"matchID"`
}

type UserSwipeNotMatchedResonse struct {
	Matched bool `json:"matched"`
}

func UserSwipe(w http.ResponseWriter, r *http.Request) {

	// Retrieve user from context
	// The AuthMiddleware is handling errors related to not finding the user
	contextUser, _ := r.Context().Value(core.UserContextKey).(models.User)

	// Only allow HTTP POST Method
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed: %s", r.Method)))
		return
	}

	var swipePayload struct {
		TargetID  uint64 `json:"targetID"`
		SwipeType string `json:"swipeType"`
	}

	err := json.NewDecoder(r.Body).Decode(&swipePayload)
	if err != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusBadRequest, fmt.Sprintf("Error decoding request body: %v", err)))
		return
	}

	// Check if the tragetID is the same as swiperID
	if swipePayload.TargetID == contextUser.ID {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusBadRequest, "Cannot swipe on yourself"))
		return
	}

	// Check if the tragetID user exists
	var targetUser models.User
	err = core.GetDb().First(&targetUser, swipePayload.TargetID).Error
	if err != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusNotFound, "Target user not found"))
		return
	}

	swipe := models.Swipe{
		SwiperID:  contextUser.ID,
		TargetID:  swipePayload.TargetID,
		SwipeType: swipePayload.SwipeType,
	}
	createSwipeErr := core.GetDb().Create(&swipe)
	if createSwipeErr.Error != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusInternalServerError, "Error creating swipe record"))
		return
	}

	// Update the Target user attractiveness score
	// Using Goroutines to run the function in the background as it's not critical to the response
	go func(targetUser models.User, swipeType string) {
		if err := updateTargetUserAttractivenessScore(&targetUser, swipeType); err != nil {
			fmt.Printf("Error updating target user Attractiveness score: %v\n", err)
		}
	}(targetUser, swipe.SwipeType)

	if swipe.SwipeType == "YES" {

		var targetSwipe models.Swipe
		err := core.GetDb().Where("swiper_id = ? AND target_id = ? AND swipe_type = ?", swipe.TargetID, swipe.SwiperID, "YES").First(&targetSwipe).Error

		// If both users swipped YES on each other
		if err == nil {
			// Check if a match already exists between these users
			if !checkIfMatchExists(swipe.SwiperID, swipe.TargetID) {
				match := models.Match{
					User1ID: swipe.SwiperID,
					User2ID: swipe.TargetID,
				}
				core.GetDb().Create(&match)
				userSwipeMatchedResonse := UserSwipeMatchedResonse{Matched: true, MatchID: match.ID}
				utils.WriteSuccessResponse(w, http.StatusOK, userSwipeMatchedResonse)
				return
			} else {
				utils.WriteErrorResponse(w, utils.NewAppError(http.StatusBadRequest, "Match already exists"))
				return
			}
		}
	}

	// The swiper selected NO so there's no need to check for a match with the target user
	createdUserResponse := UserSwipeNotMatchedResonse{Matched: false}
	utils.WriteSuccessResponse(w, http.StatusOK, createdUserResponse)
}

func updateTargetUserAttractivenessScore(targetUser *models.User, swipeType string) error {

	// Update the target user attractiveness score based on a simple formula
	// Attractiveness Score = Total Likes / (Total Likes + Total Dislikes)
	// The more likes the user gets, the higher the score
	// This is a simple example and can be improved with more complex algorithms
	fmt.Println("Updating target user attractiveness score")
	if swipeType == "YES" {
		targetUser.TotalLikesReceived++
	} else if swipeType == "NO" {
		targetUser.TotalDislikesReceived++
	}

	// Calculate attractiveness score
	targetUser.AttractivenessScore = calculateAttractivenessScore(targetUser.TotalLikesReceived, targetUser.TotalDislikesReceived)

	// Save the changes to the target user
	err := core.GetDb().Save(&targetUser).Error
	if err != nil {
		fmt.Println("Error saving target user changes", err)
		return err
	}

	return nil
}

func calculateAttractivenessScore(likes int, dislikes int) float64 {

	totalSwipes := likes + dislikes
	if totalSwipes == 0 {
		return 0.0
	}
	score := float64(likes) / float64(totalSwipes)

	// Round to the nearest 2 digits
	return math.Round(score*100) / 100
}

func checkIfMatchExists(user1ID, user2ID uint64) bool {
	var count int64
	core.GetDb().Model(&models.Match{}).Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1ID, user2ID, user2ID, user1ID).Count(&count)
	return count > 0
}
