package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/models"
	"muzz-dating/pkg/utils"
)

type PotentialMatchesResponse struct {
	ID                  uint64  `json:"id"`
	Name                string  `json:"name"`
	Gender              string  `json:"gender"`
	Age                 int     `json:"age"`
	DistanceFromMe      float64 `json:"distanceFromMe"`
	AttractivenessScore float64 `json:"attractivenessScore"`
}

func GetPotentialMatches(w http.ResponseWriter, r *http.Request) {

	// Set content-type to json
	w.Header().Set("Content-Type", "application/json")

	// Retrieve user from context
	// The AuthMiddleware is handling errors related to not finding the user
	contextUser, _ := r.Context().Value(core.UserContextKey).(models.User)

	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed: %s", r.Method)))
		return
	}

	// Get the query parameters
	minAge := r.URL.Query().Get("minAge")
	maxAge := r.URL.Query().Get("maxAge")
	gender := r.URL.Query().Get("gender")

	// Fetch all users from the database
	users := []models.User{}

	// Exclude the profiles the user matched or swiped on.
	// Exclude the user's own profile from coming up in the results
	excludedIDs := []uint64{contextUser.ID}
	swipedUserIDs := getSwipedUserIDs(contextUser.ID)
	excludedIDs = append(excludedIDs, swipedUserIDs...)
	query := core.GetDb().Omit("password", "email", "Token").Not("id IN (?)", excludedIDs)

	// Apply filters
	if minAge != "" {
		query = query.Where("age >= ?", minAge)
	}

	if maxAge != "" {
		query = query.Where("age <= ?", maxAge)
	}

	if gender != "" {
		query = query.Where("gender = ?", gender)
	}

	// Execute the query
	result := query.Find(&users)
	if err := result.Error; err != nil {
		utils.WriteErrorResponse(w, utils.NewAppError(http.StatusInternalServerError, "Error fetching users"))
		return
	}

	// Convert User slices to PotentialMatchesResponse slices
	// Used as a data transfer object to omit Token and Password fields
	potentialMatches := make([]PotentialMatchesResponse, len(users))
	for i, user := range users {
		// Calculate distance for each user and add to the result
		var distanceFromMe float64 = utils.CalculateDistance(contextUser.Latitude, contextUser.Longitude, user.Latitude, user.Longitude)
		potentialMatches[i] = PotentialMatchesResponse{
			ID:                  user.ID,
			Name:                user.Name,
			Gender:              user.Gender,
			Age:                 user.Age,
			DistanceFromMe:      distanceFromMe,
			AttractivenessScore: user.AttractivenessScore,
		}
	}

	// Sort users by distance
	sort.Slice(potentialMatches, func(i, j int) bool {
		if potentialMatches[i].AttractivenessScore == potentialMatches[j].AttractivenessScore {
			return potentialMatches[i].DistanceFromMe < potentialMatches[j].DistanceFromMe
		}
		return potentialMatches[i].AttractivenessScore > potentialMatches[j].AttractivenessScore
	})

	// Create a response object with a "results" key
	response := struct {
		Results []PotentialMatchesResponse `json:"results"`
	}{
		Results: potentialMatches,
	}
	utils.WriteSuccessResponse(w, http.StatusOK, response)
}

func getSwipedUserIDs(userID uint64) []uint64 {
	var swipes []models.Swipe
	core.GetDb().Where("swiper_id = ?", userID).Find(&swipes)

	swipedUserIDs := make([]uint64, 0, len(swipes))
	for _, swipe := range swipes {
		swipedUserIDs = append(swipedUserIDs, swipe.TargetID)
	}

	return swipedUserIDs
}
