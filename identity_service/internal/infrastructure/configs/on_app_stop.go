package configs

import (
	"context"
	"elex_storage/pkg/message_broker"

	"go.uber.org/fx"
)

func OnAppStop(lc fx.Lifecycle, rb message_broker.EventMessaging) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			rb.Close()
			return nil
		},
	})
}
