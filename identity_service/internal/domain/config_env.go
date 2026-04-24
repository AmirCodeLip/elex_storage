package domain

import "time"

type ConfigEnv struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	MigrationsDir        string
	GrpcPort             string
}
