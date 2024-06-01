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
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusBadRequest, "Payload is not valid"))
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
			userSwipeMatchedResonse := UserSwipeMatchedResonse{Matched: true, MatchID: match.ID}
			utils.WriteSuccessResponse(w, http.StatusOK, userSwipeMatchedResonse)
			return
		}
	} else {
		// Increase the Target User Total Dislikes
		targetUser.TotalDislikesReceived++
		targetUser.AttractivenessScore = calculateAttractivenessScore(targetUser.TotalLikesReceived, targetUser.TotalDislikesReceived)
		// Save the changes to the target user
		core.GetDb().Save(&targetUser)
	}

	createdUserResponse := UserSwipeNotMatchedResonse{Matched: false}
	utils.WriteSuccessResponse(w, http.StatusOK, createdUserResponse)
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
