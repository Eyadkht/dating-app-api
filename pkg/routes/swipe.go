package routes

import (
	"net/http"

	"dating-app/pkg/core"
	"dating-app/pkg/handlers"
)

func RegisterSwipeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/swipe", core.AuthMiddleware(handlers.UserSwipe))
}
