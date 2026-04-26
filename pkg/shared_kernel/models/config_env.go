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
}
