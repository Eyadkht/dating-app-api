package core

import (
	"context"
	"net/http"
	"strings"

	"dating-app/pkg/models"
	"dating-app/pkg/utils"

	"gorm.io/gorm"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteErrorResponse(w, utils.NewAppError(http.StatusUnauthorized, "Missing Authorization header"))
			return
		}

		// Split the header to get the token (assuming format: "Token <token>")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Token" {
			utils.WriteErrorResponse(w, utils.NewAppError(http.StatusUnauthorized, "Invalid Authorization header format"))
			return
		}
		token := parts[1]

		// Validate the token and bring the user object and
		// inject it into the request context so it can be used inside the protected handler
		isValid, user, err := validateToken(token)
		if err != nil {
			utils.WriteErrorResponse(w, utils.NewAppError(http.StatusInternalServerError, "Error validating token"))
			return
		}

		if !isValid {
			// Invalid token, return unauthorized
			utils.WriteErrorResponse(w, utils.NewAppError(http.StatusUnauthorized, "Invalid token"))
			return
		}

		// Token is valid, inject user info into context
		ctx := context.WithValue(r.Context(), UserContextKey, user)

		// Call the next handler with user context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func validateToken(tokenValue string) (bool, models.User, error) {

	var user models.User
	var token models.Token

	tokenResult := GetDb().Where("Value = ?", tokenValue).First(&token)
	if tokenResult.Error == gorm.ErrRecordNotFound {
		return false, user, nil
	}
	if tokenResult.Error != nil {
		return false, user, tokenResult.Error
	}

	userResult := GetDb().Where("ID = ?", token.UserID).Omit("password", "Token").First(&user)
	if userResult.Error == gorm.ErrRecordNotFound {
		return false, user, nil
	}
	if userResult.Error != nil {
		return false, user, userResult.Error
	}

	return true, user, nil
}
