package utils

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	if appErr, ok := err.(*AppError); ok {
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": appErr})
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": AppError{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error"}})
	}
}
