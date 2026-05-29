package shared_kernel

import (
	"elex_storage/pkg/shared_kernel/guard"
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
	"github.com/spf13/viper"
)

func setIdentityService(config *models.ConfigEnv) error {
	grpcUrl, err := ParseUrl(os.Getenv("IDENTITY_SERVICE_Grpc_Url"))
	config.IdentityServiceGrpcUrl = *grpcUrl
	if err != nil {
		return errors.Join(errors.New(
			fmt.Sprintf("Can't parse IDENTITY_SERVICE_Grpc_Addr is not valid host %s", grpcUrl.FullAddress)), err)
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
	apiGatewayHttpUrl, err := ParseUrl(os.Getenv("API_GATEWAY_HTTP_URL"))
	config.ApiGatewayHttpUrl = *apiGatewayHttpUrl
	if err != nil {
		return errors.Join(errors.New(
			fmt.Sprintf("Can't parse ApiGatewayServiceAddr is not valid host %s", apiGatewayHttpUrl.FullAddress)), err)
	}
	return nil
}

func setFileStorageService(config *models.ConfigEnv) error {
	fileStorageHttpUrl, err := ParseUrl(os.Getenv("FILE_STORAGE_HTTP_URL"))
	config.FileStorageHttpUrl = *fileStorageHttpUrl
	if err != nil {
		return errors.Join(errors.New(
			fmt.Sprintf("Can't parse FileMetadataServiceAddr is not valid host %s", fileStorageHttpUrl.FullAddress)), err)
	}
	return nil
}

func setFileMetadataService(config *models.ConfigEnv) error {

	// Clean DriveDisk to support linux
	parts := strings.Split(os.Getenv("DRIVE_DISK"), "\\")
	config.DriveDisk = ""
	for _, p := range parts {
		config.DriveDisk = filepath.Join(config.DriveDisk, p)
	}
	config.DriveName = os.Getenv("DRIVE_NAME")
	fileMetadataGrpcUrl, err := ParseUrl(os.Getenv("FILE_META_DATA_GRPC_URL"))
	config.FileMetadataGrpcUrl = *fileMetadataGrpcUrl
	if err != nil {
		return errors.Join(errors.New(
			fmt.Sprintf("Can't parse ApiGatewayServiceAddr is not valid host %s", fileMetadataGrpcUrl.FullAddress)), err)
	}
	return nil
}

// Load configurations from .env file.
func NewConfigEnv() (*models.ConfigEnv, error) {
	// Load config file
	config := models.ConfigEnv{}
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	configDir := filepath.Clean(filepath.Join(exeDir, "..", ".env"))
	_, err := os.Stat(configDir)
	if err == nil {
		errEnv := godotenv.Load(configDir)
		if errEnv != nil {
			err := fmt.Errorf("No .env file found in %s \n", configDir)
			log.Fatal(err)
		}
	}
	// Set services configs
	err = setIdentityService(&config)
	if err != nil {
		return nil, err
	}
	err = setApiGatewayService(&config)
	if err != nil {
		return nil, err
	}
	err = setFileMetadataService(&config)
	if err != nil {
		return nil, err
	}
	err = setFileStorageService(&config)
	if err != nil {
		return nil, err
	}
	config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")
	config.LoggerPath = os.Getenv("LOGGER_PATH")

	// Postgres configs
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

func NewConfigEnv2() (*models.ConfigEnv2, error) {
	env, err := GetPath(".env")
	if err == nil {
		errEnv := godotenv.Load(env)
		if errEnv != nil {
			err := fmt.Errorf("No .env file found in %s \n", env)
			log.Fatal(err)
		}
	}

	/// Load yml file
	ymlConfig, err := GetPath("configs.yml")
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(ymlConfig)
	viper.SetConfigType("yaml")

	// Enable environment variable expansion
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in the config
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		// Set from .env file if started with $
		if strings.Contains(val, "${") {
			expandedVal := os.ExpandEnv(val)
			viper.Set(key, expandedVal)
		}
	}

	var config2 models.ConfigEnv2
	if err := viper.Unmarshal(&config2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if migrationDir, err := GetPath("migrations"); err == nil {
		config2.MigrationsDir = migrationDir
	}

	return &config2, nil
}

func GetPath(fileName string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Start from the root
	root := filepath.VolumeName(pwd) + string(filepath.Separator)

	// Get the relative path from root to current directory
	relPath, err := filepath.Rel(root, pwd)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	// Split the relative path
	parts := strings.Split(relPath, string(filepath.Separator))

	// Find where "cmd" appears and truncate
	for i, part := range parts {
		if part == "cmd" {
			parts = parts[:i]
			break
		}
	}

	// Build the final path
	filePath := filepath.Join(root, filepath.Join(parts...), fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file does not exist: %s", filePath)
		}
		return "", fmt.Errorf("failed to check file existence: %w", err)
	}

	return filePath, nil
}

func GetTestEnvPath() (*string, error) {
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	paths := strings.Split(exeDir, "\\")
	for i := 0; ; i++ {
		/// remove until get root of project
		last := paths[len(paths)-1]
		if last != "elex_storage" {
			paths = paths[:len(paths)-1]
		} else {
			break
		}
		if i == 7 {
			return nil, errors.New("Can't find env path")
		}
	}
	paths = append(paths, ".env")
	envAddr := strings.Join(paths, "\\")
	return &envAddr, nil
}

func TestConfigEnv(envFilePath *string) (*models.ConfigEnv, error) {
	if guard.AgainstPNullStr(envFilePath) {
		if err := godotenv.Load(*envFilePath); err != nil {
			return nil, errors.New(fmt.Sprintf("Warning: Could not load %s", *envFilePath))
		}
	} else {
		return nil, errors.New("Provide valid .env path")
	}

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

	return &config, nil
}

func ParseUrl(addr string) (*models.Url, error) {
	if !guard.AgainstEmptyStr(addr) {
		return nil, errors.New("Addr can't be nil")
	}
	result := models.Url{
		Protocol: "http", // default protocol
		Port:     80,     // default port
	}

	// Check for protocol
	if strings.Contains(addr, "://") {
		parts := strings.SplitN(addr, "://", 2)
		result.Protocol = parts[0]
		addr = parts[1]
	}

	// Remove trailing slash
	addr = strings.TrimSuffix(addr, "/")

	// Parse host and port
	if strings.Contains(addr, ":") {
		parts := strings.Split(addr, ":")
		result.Host = parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return &models.Url{}, err
		}
		result.Port = port
	} else {
		// Check if it's just a port number
		if port, err := strconv.Atoi(addr); err == nil {
			result.Port = port
			result.Host = "localhost"
		} else {
			// It's just a host
			result.Host = addr
		}
	}

	// Build full URL
	if result.Port == 80 && result.Protocol == "http" {
		result.FullAddress = fmt.Sprintf("%s://%s", result.Protocol, result.Host)
		result.Address = result.Host
	} else if result.Port == 443 && result.Protocol == "https" {
		result.FullAddress = fmt.Sprintf("%s://%s", result.Protocol, result.Host)
		result.Address = fmt.Sprintf("%s", result.Host)
	} else {
		result.FullAddress = fmt.Sprintf("%s://%s:%d", result.Protocol, result.Host, result.Port)
		result.Address = fmt.Sprintf("%s:%d", result.Host, result.Port)
	}

	return &result, nil
}
