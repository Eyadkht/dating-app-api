package routes

import (
	"muzz-dating/pkg/handlers"
	"net/http"
)

func RegisterDiscoverRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/discover", handlers.GetPotentialMatches)
}
