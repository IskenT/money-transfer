package config

import (
	"os"
	"strconv"
	"time"
)

// Config
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig
type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// NewConfig
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     getEnvAsDuration("READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "money_transfer"),
			Password: getEnv("POSTGRES_PASSWORD", "password"),
			Name:     getEnv("POSTGRES_DB", "money_transfer"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// getEnv
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsBool
func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return b
	}
	return fallback
}

// getEnvAsDuration
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		d, err := time.ParseDuration(value)
		if err != nil {
			return fallback
		}
		return d
	}
	return fallback
}
