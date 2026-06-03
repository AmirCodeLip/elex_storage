package use_case

import (
	"elex_storage/identity_service/internal/use_case/cqrs/command_handlers"

	"go.uber.org/fx"
)

func UseCaseModule() fx.Option {
	var module = fx.Options(
		fx.Provide(command_handlers.NewRegisterUserHandler),
		fx.Provide(command_handlers.NewLoginUserHandler),
	)
	return module
}
