package database

import (
	"fmt"
	"log"
	"time"

	"github.com/IskenT/money-transfer/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// DBConfig
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDBConfig
func NewDBConfig(cfg *config.Config) *DBConfig {
	return &DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}
}

// DSN
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// NewDB
func NewDB(dbConfig *DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dbConfig.DSN())
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// NewDBWithRetry
func NewDBWithRetry(dbConfig *DBConfig, maxRetries int, retryInterval time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = NewDB(dbConfig)
		if err == nil {
			return db, nil
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
