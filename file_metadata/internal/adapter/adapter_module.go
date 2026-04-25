package adapter

import (
	"elex_storage/file_metadata/internal/adapter/grpc_server"
	"elex_storage/file_metadata/internal/adapter/messaging/subscribers"
	"elex_storage/file_metadata/internal/adapter/pgx_repositories"
	"elex_storage/pkg/message_broker"
	"net/http"

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

func registerSubscribes(eventMessaging message_broker.EventMessaging, fileCreatedSubscriber *subscribers.FileCreatedSubscriber) {
	eventMessaging.Subscribe(nil, "file.created", fileCreatedSubscriber.Handle)
}

func AdapterModule() fx.Option {
	var module = fx.Options(
		// Provide repositores
		fx.Provide(pgx_repositories.CreateFileMetadataRepository),
		fx.Provide(pgx_repositories.CreateDirectoryMetadataRepository),
		// Provide handlers
		fx.Provide(grpc_server.NewFileMetadataService),
		// Messaging setup
		fx.Provide(subscribers.NewFileCreatedSubscriber),
		fx.Invoke(registerSubscribes),
	)
	return module
}

// Middleware example
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
