package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct is a type that server is using for its configuration.
type Config struct {
	ServerPort            string
	JWTsecretKey          string
	UploadsLocation       string
	MigrationPath         string
	MetadataPQCredentials PostgresCredentials
	UserPQCredentials     PostgresCredentials
}

type PQStorage string

const (
	MetadataPQStorage PQStorage = "PQ_METADATADB"
	UserPQStorage     PQStorage = "PQ_USERDB"
)

// LoadConfig loads config from .env file.
func (c *Config) LoadConfig(configPath string) {
	if configPath == "" {
		log.Fatal("You must provide --config with the path to the config file")
	}

	// load config file
	if err := godotenv.Load(configPath); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	c.ServerPort = getEnv("SERVER_PORT", ":3000")
	c.JWTsecretKey = getEnv("SECRET_KEY", "als;dgasdfkasbf2ql4q")
	c.UploadsLocation = getEnv("DISK_UPLOADS_LOCATION", "uploads")
	c.MigrationPath = getEnv("MIGRATION_PATH", "migrations")
	c.UserPQCredentials = getPostgresCredentials(UserPQStorage)
	c.MetadataPQCredentials = getPostgresCredentials(MetadataPQStorage)
}

// getEnv set key value from env if exists, otherwise default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getPostgresCredentials returns credentials that are used to connect to db.
func getPostgresCredentials(dbName PQStorage) PostgresCredentials {
	return PostgresCredentials{
		DB_USER:     getEnv(fmt.Sprintf("%s_USER", dbName), "admin"),
		DB_PASSWORD: getEnv(fmt.Sprintf("%s_PASSWORD", dbName), "root"),
		DB_NAME:     getEnv(fmt.Sprintf("%s_NAME", dbName), "fyleDB"),
		DB_HOST:     getEnv(fmt.Sprintf("%s_HOST", dbName), "postgres"),
		DB_PORT:     getEnv(fmt.Sprintf("%s_PORT", dbName), "5432"),
	}
}
