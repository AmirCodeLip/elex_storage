package main

import (
	"elex_storage/identity_service/internal/core_utils"
	"elex_storage/identity_service/internal/domain"
	"elex_storage/identity_service/internal/infrastructure/configs"
	"elex_storage/identity_service/internal/infrastructure/database"
	"fmt"

	"elex_storage/identity_service/internal/adapter"
	"elex_storage/identity_service/internal/use_case"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/message_broker"
	"elex_storage/pkg/shared_kernel/utils"

	"go.uber.org/fx"
)

func main() {
	// before start set .env file variables
	cfg, err := configs.NewConfigEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	fx.New(
		fx.Provide(
			func() *domain.ConfigEnv {
				return cfg
			},
			logger.NewLogger,
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
