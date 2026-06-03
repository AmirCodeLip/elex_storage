package main

import (
	"elex_storage/api_gateway/internal/adapter"
	"elex_storage/api_gateway/internal/infrastructure/configs"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/utils"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			shared_kernel.NewConfigEnv,
			logger.NewLogger,
			configs.NewServer,
			configs.NewGRPCClient,
			utils.NewHttpErrorUtils,
		),
		adapter.AdapterModule(),
		fx.Invoke(configs.RegisterFX),
	).Run()
}
