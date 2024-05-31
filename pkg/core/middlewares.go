package core

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"muzz-dating/pkg/models"

	"gorm.io/gorm"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set content-type to json
		w.Header().Set("Content-Type", "application/json")

		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Missing token, return unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing Authorization header"})
			return
		}

		// Split the header to get the token (assuming format: "Token <token>")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Token" {
			// Invalid format, return unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Authorization header format"})
			return
		}
		token := parts[1]

		// Validate the token (replace with your token validation logic)
		// This example is a placeholder, you'll need to implement actual validation
		isValid, user, err := validateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error validating token"})
			return
		}

		if !isValid {
			// Invalid token, return unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
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
