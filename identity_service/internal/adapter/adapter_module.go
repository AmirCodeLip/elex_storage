package adapter

import (
	"elex_storage/identity_service/internal/adapter/grpc_server"
	"elex_storage/identity_service/internal/adapter/messaging/publishers"
	"elex_storage/identity_service/internal/adapter/pgx_repositories"
	"elex_storage/pkg/shared_kernel/guard"
	"elex_storage/pkg/shared_kernel/models"
	"elex_storage/pkg/shared_kernel/token_handlers"
	"errors"

	"go.uber.org/fx"
)

func configureJwt(config *models.ConfigEnv) (*token_handlers.JwtToken, error) {
	if !guard.AgainstTimeDurationPtr(config.AccessTokenDuration) || !guard.AgainstTimeDurationPtr(config.RefreshTokenDuration) {
		return nil, errors.New(`The AccessTokenDuration or RefreshTokenDuration is either missing or invalid.
		 Please configure them in config.yml`)
	}
	jwtToken, err := token_handlers.NewJwtToken(*config.AccessTokenDuration, *config.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}
	if !guard.AgainstPNullStr(config.JWTPrivateKeyPath) || !guard.AgainstPNullStr(config.JWTPublicKeyPath) {
		return nil, errors.New(`The JWTPrivateKeyPath or JWTPublicKeyPath is either missing or invalid.
		 Please configure them in config.yml`)
	}
	if err != nil {
		return nil, err
	}
	err = jwtToken.SetPrivateKey(*config.JWTPrivateKeyPath)
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
