package models

import (
	"time"
)

type ConfigEnv struct {
	MigrationsDir            string
	LoggerPath               string
	PGConnectionString       string
	RabbitmqConnectionString string
	// identity_service configs
	// IdentityServiceGrpcHost string
	// IdentityServiceGrpcPort string

	IdentityServiceGrpcUrl Url
	AccessTokenDuration    time.Duration
	RefreshTokenDuration   time.Duration
	// api_gateway configs
	ApiGatewayHttpUrl Url
	//  file_metadata configs
	FileMetadataGrpcUrl Url
	// file_storage configs
	DriveDisk          string
	DriveName          string
	FileStorageHttpUrl Url
	// Monitoring and logger configs.
	LokiApiAddress string
}

type ConfigEnv2 struct {
	ServiceName   string `mapstructure:"service_name"`
	DriveDisk     string `mapstructure:"drive_disk"`
	MigrationsDir string
	Loki          Loki     `mapstructure:"loki"`
	Server        Server   `mapstructure:"server"`
	Database      Database `mapstructure:"database"`
	RabbitMQ      RabbitMQ `mapstructure:"rabbitmq"`
}
