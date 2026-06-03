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

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func NewConfigEnv() (*models.ConfigEnv, error) {
	env, err := GetRelativePath(".env")
	if err == nil {
		errEnv := godotenv.Load(env)
		if errEnv != nil {
			err := fmt.Errorf("No .env file found in %s \n", env)
			log.Fatal(err)
		}
	}

	/// Load yml file
	ymlConfig, err := GetRelativePath("configs.yml")
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

	var config models.ConfigEnv
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if migrationDir, err := GetRelativePath("migrations"); err == nil {
		config.MigrationsDir = migrationDir
	}

	return &config, nil
}

// This function will return file path witout cmd file
func GetRelativePath(fileName string) (string, error) {
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
	return nil, errors.New("Not implemented new config")
	// config := models.ConfigEnv{}
	// os.Setenv("DB_HOST", "localhost")
	// os.Setenv("DB_PORT", "10252")
	// os.Setenv("DB_DATABASE", "file_metadata")
	// os.Setenv("DB_USERNAME", "elex_storage")
	// os.Setenv("DB_PASSWORD", "pass1234")
	// os.Setenv("DB_SCHEMA", "public")
	// os.Setenv("MIGRATIONS_DIR", "..\\migrations")

	// config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")
	// database := os.Getenv("DB_DATABASE")
	// password := os.Getenv("DB_PASSWORD")
	// username := os.Getenv("DB_USERNAME")
	// port := os.Getenv("DB_PORT")
	// host := os.Getenv("DB_HOST")
	// schema := os.Getenv("DB_SCHEMA")
	// pgConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	// config.PGConnectionString = pgConnectionString
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
