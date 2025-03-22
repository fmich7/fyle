package config

import (
	"log"
	"os"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/joho/godotenv"
)

// Config struct is a type that server is using for its configuration
type Config struct {
	ServerPort          string
	JWTsecretKey        string
	UploadsLocation     string
	MigrationPath       string
	PostgresCredentials types.PostgresCredentials
}

// LoadConfig loads config from .env file
func (c *Config) LoadConfig(fileName string) {
	if fileName == "" {
		fileName = ".env"
	}

	// load config file
	if err := godotenv.Load(fileName); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	c.ServerPort = getEnv("SERVER_PORT", ":3000")
	c.JWTsecretKey = getEnv("SECRET_KEY", "als;dgasdfkasbf2ql4q")
	c.UploadsLocation = getEnv("DISK_UPLOADS_LOCATION", "uploads")
	c.MigrationPath = getEnv("MIGRATION_PATH", "migrations")
	c.PostgresCredentials = getPostgresCredentials()
}

// getEnv set key value from env if exists, otherwise default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getPostgresCredentials returns credentials that are used to connect to db
func getPostgresCredentials() types.PostgresCredentials {
	return types.PostgresCredentials{
		DB_USER:     getEnv("POSTGRES_USER", "admin"),
		DB_PASSWORD: getEnv("POSTGRES_PASSWORD", "root"),
		DB_NAME:     getEnv("POSTGRES_NAME", "fyleDB"),
		DB_HOST:     getEnv("POSTGRES_HOST", "postgres"),
		DB_PORT:     getEnv("POSTGRES_PORT", "5432"),
	}
}
