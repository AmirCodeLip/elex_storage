package models

import "time"

type ConfigEnv struct {
	MigrationsDir            string
	LoggerPath               string
	PGConnectionString       string
	RabbitmqConnectionString string

	// identity_service configs
	IdentityServiceGrpcHost string
	IdentityServiceGrpcPort string
	IdentityServiceGrpcAddr string
	AccessTokenDuration     time.Duration
	RefreshTokenDuration    time.Duration
	// api_gateway configs
	ApiGatewayHttpAddr string
	ApiGatewayHttpPort string
	ApiGatewayHttpHost string
	//  file_metadata configs
	FileMetadataGrpcAddr string
	FileMetadataGrpcHost string
	FileMetadataGrpcPort string
	// file_storage configs
	DriveDisk           string
	DriveName           string
	FileStorageHttpAddr string
	FileStorageHttpHost string
	FileStorageHttpPort string
}
