package main

import (
	"elex_storage/file_storage/internal/core_utils"
	"elex_storage/file_storage/internal/infrastructure/configs"
	"elex_storage/file_storage/internal/infrastructure/database"

	"elex_storage/file_storage/internal/adapter"
	"elex_storage/file_storage/internal/use_case"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/message_broker"

	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/utils"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			shared_kernel.NewConfigEnv2,
			logger.NewLokiLogger,
			database.NewDatabase,
			core_utils.NewPathUtil,
			configs.NewServer,
			message_broker.NewRabbitmqMessaging,
			utils.NewHttpErrorUtils,
		),
		use_case.UseCaseModule(),
		adapter.AdapterModule(),
		fx.Invoke(database.ApplyMigration),
		fx.Invoke(configs.RegisterFX),
		fx.Invoke(configs.OnAppStop),
	).Run()
}
