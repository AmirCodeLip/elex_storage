package adapter

import (
	"context"
	"elex_storage/api_gateway/internal/adapter/http_clients"
	"elex_storage/api_gateway/internal/adapter/http_handlers"
	"elex_storage/api_gateway/internal/domain/client_repositories"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/fx"
	"google.golang.org/grpc"
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

type GrpcInput[Req any, Res any] func(
	ctx context.Context,
	in *Req,
	opts ...grpc.CallOption,
) (*Res, error)

func RegisterHttpRoutes(httpHandler *http_handlers.HttpHandler, fileStorageRepository client_repositories.FileStorageRepository) *http.Handler {
	r := httprouter.New()
	corsWrapper := corsMiddleware(r)

	// Register routes
	// User routes
	r.HandlerFunc(http.MethodPost, "/users/login", http_handlers.WrapGrpcFunc(httpHandler, httpHandler.UserServiceClient.Login))
	r.HandlerFunc(http.MethodPost, "/users/register", http_handlers.WrapGrpcFunc(httpHandler, httpHandler.UserServiceClient.Register))
	// File storage routes
	r.HandlerFunc(http.MethodPost, "/file/upload", http_handlers.HttpProxyFunc(httpHandler, fileStorageRepository.Upload))

	return &corsWrapper
}

func AdapterModule() fx.Option {

	var module = fx.Options(
		// Provide http handlers
		fx.Provide(http_clients.NewFileStorageClient),
		fx.Provide(http_handlers.NewHttpHandler),
		fx.Provide(RegisterHttpRoutes),
	)
	return module
}
