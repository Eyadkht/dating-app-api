package utils

import (
	"encoding/json"
	"net/http"
)

func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		WriteErrorResponse(w, NewAppError(http.StatusInternalServerError, "Error serializing response to JSON"))
	}
}
