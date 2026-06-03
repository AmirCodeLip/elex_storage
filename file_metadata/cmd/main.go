package main

import (
	"elex_storage/file_metadata/internal/adapter"
	"elex_storage/file_metadata/internal/infrastructure/configs"
	"elex_storage/file_metadata/internal/infrastructure/database"
	"elex_storage/file_metadata/internal/use_case"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/message_broker"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/utils"

	"go.uber.org/fx"
)

func main() {
	// before start set .env file variables
	fx.New(
		fx.Provide(
			shared_kernel.NewConfigEnv,
			logger.NewLokiLogger,
			database.NewDatabase,
			utils.NewHttpErrorUtils,
			message_broker.NewRabbitmqMessaging,
			configs.NewGRPCServer,
		),
		use_case.UseCaseModule(),
		adapter.AdapterModule(),
		fx.Invoke(database.ApplyMigration),
		fx.Invoke(configs.RegisterAndStartGRPCServer),
		fx.Invoke(configs.OnAppStop),
	).Run()
}
