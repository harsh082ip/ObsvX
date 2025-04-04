package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	LogLevel     string
	LogFile      string
	LogToConsole bool
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using defaults")
	}

	config := &AppConfig{
		DBHost:       getEnv("DB_HOST", "host"),
		DBPort:       getEnv("DB_PORT", "5000"),
		DBUser:       getEnv("DB_USER", "admin"),
		DBPassword:   getEnv("DB_PASSWORD", "pass"),
		DBName:       getEnv("DB_NAME", "default"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		LogFile:      getEnv("LOG_FILE", "service_logs.log"),
		LogToConsole: getEnvBool("LOG_TO_CONSOLE", true),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1" || value == "yes"
	}
	return fallback
}
