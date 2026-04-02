package environment

import (
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
    EnvKey         = "ENV" 
    ServicePortKey = "SERVICE_PORT"
    VersionKey     = "VERSION"

	// Database keys 
	DbHostKey      = "DB_HOST"
	DbPortKey      = "DB_PORT"
	DbNameKey      = "DB_NAME"
	DbUserKey      = "DB_USER"
	DbPassKey      = "DB_PASS"

    // Login token
	LoginAccessExpMinsKey            = "LOGIN_ACCESS_EXP_MINS"
	LoginRefreshExpMinsKey           = "LOGIN_REFRESH_EXP_MINS"
	LoginRefreshExpMinsRememberMeKey = "LOGIN_REFRESH_EXP_MINS_REMEMBER_ME"
	LoginAccessSecretKey             = "LOGIN_ACCESS_SECRET"
	LoginRefreshSecretKey            = "LOGIN_REFRESH_SECRET"
    
    // CORS
    AllowOriginKey     = "ALLOW_ORIGINS"
    AllowCredentialKey = "ALLOW_CREDENTIALS"

    RequestMaxBodySizeMB = "REQUEST_MAX_BODY_SIZE_MB"
)

func New(dirDepth uint) {
	// Enable automatic environment variable reading first
	// This ensures env vars are available as fallback
	viper.AutomaticEnv()

	// Try to read from .env file
	viper.SetConfigFile(".env")

	var configDir string
	if dirDepth == 0 {
		configDir = "."
	} else {
		configDir = ".."
		for i := uint(1); i < dirDepth; i++ {
			configDir = filepath.Join(configDir, "..")
		}
	}

	viper.AddConfigPath(configDir)

	// Try to read config file, but don't fail if not found
	if err := viper.ReadInConfig(); err != nil {
		// .env file not found or error reading, use environment variables
		log.Println(".env file not found, using system environment variables")
	} else {
		log.Println("Loaded config from file:", viper.ConfigFileUsed())
	}
}

func GetString(key string) string {
	if !viper.IsSet(key) {
		panic("Failed to get environment key: " + key)
	}

	return viper.GetString(key)
}

func GetInt(key string) int {
	if !viper.IsSet(key) {
		panic("Failed to get environment key: " + key)
	}

	return viper.GetInt(key)
}

func GetBool(key string) bool {
	if !viper.IsSet(key) {
		panic("Failed to get environment key: " + key)
	}

	return viper.GetBool(key)
}

func GetUuid(key string) uuid.UUID {
	val := GetString(key)
	id, err := uuid.Parse(val)
	if err != nil {
		panic("Failed to parse environment key as UUID: " + key + " (" + val + ")")
	}
	return id
}

func GetIntWithoutPanic(key string) int {
	if !viper.IsSet(key) {
		return 0
	}

	return viper.GetInt(key)
}

func GetRequestMaxBodySizeLimit(key string) int {
	const megabyteUnit = 1024 * 1024
	const defaultMB = 25
	const minMB = 10

	if !viper.IsSet(key) {
		return defaultMB * megabyteUnit
	}

	mb := viper.GetInt(key)

	if mb < minMB {
		return minMB * megabyteUnit
	}

	return mb * megabyteUnit
}

func GetVersion() string {
	if viper.IsSet(VersionKey) {
		return viper.GetString(VersionKey)
	}
	return "1.0.0-fallback"
}
