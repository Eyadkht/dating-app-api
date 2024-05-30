package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "Hello World!"}
		json.NewEncoder(w).Encode(response)
	})

	//Run Server
	fmt.Println("Server is running on port 8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}