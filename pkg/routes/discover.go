package routes

import (
	"net/http"

	"dating-app/pkg/core"
	"dating-app/pkg/handlers"
)

func RegisterDiscoverRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/discover", core.AuthMiddleware(handlers.GetPotentialMatches))
}
