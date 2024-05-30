package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: CreateUser")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed: %s", r.Method)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Hello World!"}
	json.NewEncoder(w).Encode(response)
}
