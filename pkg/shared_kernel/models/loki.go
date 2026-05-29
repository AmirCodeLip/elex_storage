package models

type Loki struct {
	APIAddress   string `mapstructure:"api_address"`
	QueryAddress string `mapstructure:"query_address"`
	Auth         Auth   `mapstructure:"auth"`
}

type Auth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	TenantID string `mapstructure:"tenant_id"`
}

type Server struct {
	HTTPListenUrl string `mapstructure:"http_listen_url"`
	GRPCListenUrl string `mapstructure:"grpc_listen_url"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Schema   string `mapstructure:"schema"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
