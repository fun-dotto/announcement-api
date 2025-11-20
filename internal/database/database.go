package database

import (
	"fmt"
	"log"

	"github.com/fun-dotto/announcement-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("Connected to database: %s@%s:%s/%s", cfg.User, cfg.Host, cfg.Port, cfg.DBName)
	return db, nil
}

// データベース接続を閉じる
func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
