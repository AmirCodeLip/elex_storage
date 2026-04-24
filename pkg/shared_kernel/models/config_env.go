package models

type ConfigEnv struct {
	MigrationsDir            string
	GrpcPort                 string
	HttpPort                 string
	LoggerPath               string
	PGConnectionString       string
	RabbitmqConnectionString string

	/// identity_service configs
	IdentityServiceGrpcHost string
	IdentityServiceGrpcPort string
	IdentityServiceGrpcAddr string

	/// api_gateway configs
	ApiGatewayServiceAddr string
	ApiGatewayServicePort string
	ApiGatewayServiceHost string
	//  Http_addr
}
