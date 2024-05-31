package main

import (
	"fmt"
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/routes"
)

func main() {

	// Load Environment variables
	core.LoadConfig()

	// Initiate Db Connection
	fmt.Println("Establishing Database connection")
	core.InitDb()

	// Initiate Routers
	fmt.Println("Registering Routes")
	mux := http.NewServeMux()
	routes.RegisterUserRoutes(mux)
	routes.RegisterDiscoverRoutes(mux)

	// Run Server
	fmt.Println("Server is running on port 8888")
	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
