package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct is a type that server is using for its configuration
type Config struct {
	ServerPort      string
	JWTsecretKey    string
	UploadsLocation string
}

// LoadConfig loads config to Config from .env file with specified fileName
func (c *Config) LoadConfig(fileName string) {
	if fileName == "" {
		fileName = ".env"
	}

	// load config
	if err := godotenv.Load(fileName); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	c.ServerPort = getEnv("SERVER_PORT", ":8080")
	c.JWTsecretKey = getEnv("SECRET_KEY", "als;dgasdfkasbf2ql4q")
	c.UploadsLocation = getEnv("DISK_UPLOADS_LOCATION", "uploads")
}

// getEnv looks up if key exists in env, if so returns it;s value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// NewTestingConfig return config that is used for testing
func NewTestingConfig() *Config {
	return &Config{
		ServerPort: ":0",
	}
}
