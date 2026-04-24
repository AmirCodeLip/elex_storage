package use_case

import (
	"elex_storage/file_metadata/internal/use_case/cqrs/command_handlers"

	"go.uber.org/fx"
)

func UseCaseModule() fx.Option {
	var module = fx.Options(
		fx.Provide(command_handlers.NewCreateDirectoryHandler),
		fx.Provide(command_handlers.NewFileCreatedHandler),
		fx.Provide(command_handlers.NewGetDirectoriesHandler),
	)
	return module
}
