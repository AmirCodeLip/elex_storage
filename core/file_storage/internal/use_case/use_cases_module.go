package use_case

import (
	"elex_storage/file_storage/internal/use_case/cqrs/command_handlers"

	"go.uber.org/fx"
)

func UseCaseModule() fx.Option {
	var module = fx.Options(
		fx.Provide(command_handlers.NewSaveFileHandler),
	)
	return module
}
