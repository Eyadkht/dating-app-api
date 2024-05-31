package routes

import (
	"net/http"

	"muzz-dating/pkg/core"
	"muzz-dating/pkg/handlers"
)

func RegisterSwipeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/swipe", core.AuthMiddleware(handlers.UserSwipe))
}
