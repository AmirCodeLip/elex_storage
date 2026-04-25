package shared_kernel

import (
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func setIdentityService(config *models.ConfigEnv) error {
	grpcAddr := os.Getenv("IDENTITY_SERVICE_Grpc_Addr")
	config.IdentityServiceGrpcAddr = grpcAddr
	host, port, err := parseAddress(grpcAddr)
	config.IdentityServiceGrpcHost = host
	config.IdentityServiceGrpcPort = fmt.Sprintf("%d", port)
	if err != nil {
		return errors.Join(errors.New(
			fmt.Sprintf("Can't parse IDENTITY_SERVICE_Grpc_Addr is not valid host %s", grpcAddr)), err)
	}
	// Set AccessTokenDuration
	val1 := os.Getenv("ACCESS_TOKEN_DURATION")
	accessTokenDuration, err := time.ParseDuration(val1)
	if err != nil {
		return err
	}
	val2 := os.Getenv("REFRESH_TOKEN_DURATION")
	refreshTokenDuration, err := time.ParseDuration(val2)
	if err != nil {
		return err
	}
	config.RefreshTokenDuration = refreshTokenDuration
	config.AccessTokenDuration = accessTokenDuration
	return nil
}

func setApiGatewayService(config *models.ConfigEnv) error {
	api_gateway_http_addr := os.Getenv("API_GATEWAY_HTTP_Addr")
	config.ApiGatewayServiceAddr = api_gateway_http_addr
	host, port, err := parseAddress(api_gateway_http_addr)
	config.ApiGatewayServiceHost = host
	config.ApiGatewayServicePort = fmt.Sprintf("%d", port)
	if err != nil {
		return errors.Join(errors.New(
			fmt.Sprintf("Can't parse ApiGatewayServiceAddr is not valid host %s", config.ApiGatewayServiceAddr)), err)
	}
	return nil
}

// Load configurations from .env file.
func NewConfigEnv() (*models.ConfigEnv, error) {
	/// Load config file
	config := models.ConfigEnv{}
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	configDir := filepath.Join(exeDir, ".env")
	errEnv := godotenv.Load(configDir)
	if errEnv != nil {
		err := fmt.Errorf("No .env file found in %s \n", configDir)
		log.Fatal(err)
	}

	/// Set services configs
	err := setIdentityService(&config)
	if err != nil {
		return nil, err
	}
	err = setApiGatewayService(&config)
	if err != nil {
		return nil, err
	}

	config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")
	config.LoggerPath = os.Getenv("LOGGER_PATH")
	config.HttpPort = os.Getenv("FILE_META_DATA_HTTP_PORT")
	config.GrpcPort = os.Getenv("FILE_META_DATA_GRPC_PORT")

	/// Postgres configs
	database := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	schema := os.Getenv("DB_SCHEMA")
	pgConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	config.PGConnectionString = pgConnectionString

	rabbitmqUserName := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	rabbitmqConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUserName, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
	config.RabbitmqConnectionString = rabbitmqConnectionString

	return &config, nil
}

func TestConfigEnv() *models.ConfigEnv {
	config := models.ConfigEnv{}
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "10252")
	os.Setenv("DB_DATABASE", "file_metadata")
	os.Setenv("DB_USERNAME", "elex_storage")
	os.Setenv("DB_PASSWORD", "pass1234")
	os.Setenv("DB_SCHEMA", "public")
	os.Setenv("MIGRATIONS_DIR", "..\\migrations")

	config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")
	database := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	schema := os.Getenv("DB_SCHEMA")
	pgConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	config.PGConnectionString = pgConnectionString

	return &config
}

func parseAddress(addr string) (host string, port int, err error) {
	if strings.Contains(addr, ":") {
		parts := strings.Split(addr, ":")
		host = parts[0]
		port, err = strconv.Atoi(parts[1])
		return
	}

	// no colon
	if p, err := strconv.Atoi(addr); err == nil {
		// it's a port
		return "", p, nil
	}

	// it's a host
	return addr, 0, nil
}
