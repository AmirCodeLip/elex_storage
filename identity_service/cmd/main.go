package main

import (
	"elex_storage/identity_service/internal/core_utils"
	"elex_storage/identity_service/internal/infrastructure/configs"
	"elex_storage/identity_service/internal/infrastructure/database"

	"elex_storage/identity_service/internal/adapter"
	"elex_storage/identity_service/internal/use_case"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/message_broker"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/utils"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			shared_kernel.NewConfigEnv,
			logger.NewLokiLogger,
			database.NewDatabase,
			utils.NewHttpErrorUtils,
			message_broker.NewRabbitmqMessaging,
			core_utils.NewIdentityUtil,
			configs.NewGRPCServer,
		),
		use_case.UseCaseModule(),
		adapter.AdapterModule(),
		fx.Invoke(database.ApplyMigration),
		fx.Invoke(configs.RegisterAndStartGRPCServer),
		fx.Invoke(configs.OnAppStop),
	).Run()
}
