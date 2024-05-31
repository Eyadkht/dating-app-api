package routes

import (
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/handlers"
)

func RegisterDiscoverRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/discover", core.AuthMiddleware(handlers.GetPotentialMatches))
}
