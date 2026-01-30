// Package routes provides centralized route registration for the MrRSS API.
// This eliminates code duplication between main.go and main-core.go.
package routes

import (
	"net/http"

	"MrRSS/internal/handlers/core"
)

// RegisterAPIRoutes registers all API routes to the provided mux.
// This function is called by both main.go (desktop mode) and main-core.go (server mode).
func RegisterAPIRoutes(mux *http.ServeMux, h *core.Handler) {
	registerFeedRoutes(mux, h)
	registerArticleRoutes(mux, h)
	registerAIRoutes(mux, h)
	registerSettingsRoutes(mux, h)
	registerOtherRoutes(mux, h)
}
