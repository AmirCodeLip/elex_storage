package adapter

import (
	"elex_storage/file_storage/internal/adapter/http_handlers"
	"elex_storage/file_storage/internal/adapter/messaging/publishers"
	"elex_storage/file_storage/internal/adapter/pgx_repositories"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/fx"
)

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Use "*" for all origins, or replace with specific origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are needed

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterHttpRoutes(httpHandler *http_handlers.HttpHandler) *http.Handler {
	r := httprouter.New()
	corsWrapper := corsMiddleware(r)
	// Register routes
	r.HandlerFunc(http.MethodPost, "/upload", httpHandler.UploadHandler)
	return &corsWrapper
}

func AdapterModule() fx.Option {
	var module = fx.Options(
		// Provide repositores
		fx.Provide(pgx_repositories.CreateFileRepository),
		// Provide http handlers
		fx.Provide(http_handlers.NewHttpHandler),
		fx.Provide(RegisterHttpRoutes),
		// Provide messaging
		fx.Provide(publishers.CreateStoragePublisher),
	)

	return module
}
