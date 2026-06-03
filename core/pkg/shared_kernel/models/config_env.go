package models

import (
	"time"
)

type ConfigEnv struct {
	ServiceName string `mapstructure:"service_name"`

	MigrationsDir string
	Loki          Loki     `mapstructure:"loki"`
	Server        Server   `mapstructure:"server"`
	Database      Database `mapstructure:"database"`
	RabbitMQ      RabbitMQ `mapstructure:"rabbitmq"`
	// file_storage configs
	DriveDisk string `mapstructure:"drive_disk"`
	// Identity configs
	AccessTokenDuration  *time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration *time.Duration `mapstructure:"refresh_token_duration"`
	JWTPrivateKeyPath    *string        `mapstructure:"jwt_private_key_path"`
	JWTPublicKeyPath     *string        `mapstructure:"jwt_public_key_path"`
}
