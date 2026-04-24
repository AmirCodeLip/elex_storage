package adapter

import (
	"elex_storage/identity_service/internal/adapter/grpc_server"
	"elex_storage/identity_service/internal/adapter/messaging/publishers"
	"elex_storage/identity_service/internal/adapter/pgx_repositories"
	"elex_storage/identity_service/internal/domain"
	"elex_storage/pkg/shared_kernel/token_handlers"
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

func configureJwt(config *domain.ConfigEnv) (*token_handlers.JwtToken, error) {
	jwtToken, err := token_handlers.NewJwtToken(config.AccessTokenDuration, config.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}
	err = jwtToken.SetPublicKey()
	if err != nil {
		return nil, err
	}
	err = jwtToken.SetPrivateKey()
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}

func AdapterModule() fx.Option {
	var module = fx.Options(
		// Provide repositores
		fx.Provide(pgx_repositories.CreateFileRepository),
		// Provide grpc handlers
		fx.Provide(grpc_server.NewUserService),
		// Provide messaging
		fx.Provide(publishers.CreateStoragePublisher),
		fx.Provide(configureJwt),
	)
	return module
}
