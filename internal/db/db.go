package db

import (
	"fmt"
	"time"

	"github.com/harsh082ip/ObsvX/config"
	"github.com/harsh082ip/ObsvX/internal/log"
	"github.com/harsh082ip/ObsvX/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDB(cfg *config.AppConfig) (*gorm.DB, error) {
	logger := log.InitLogger("database")
	start := time.Now()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	logger.LogInfoMessage().
		Str("host", cfg.DBHost).
		Str("port", cfg.DBPort).
		Str("database", cfg.DBName).
		Msg("Connecting to database")

	// Configure GORM logger
	gormLogger := gormlogger.Default.LogMode(gormlogger.Info)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		logger.LogErrorMessage().
			Str("error", err.Error()).
			Str("host", cfg.DBHost).
			Str("port", cfg.DBPort).
			Str("database", cfg.DBName).
			Dur("attempt_duration", time.Since(start)).
			Msg("Failed to connect to database")
		return nil, err
	}

	// Auto Migrate models
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		logger.LogErrorMessage().
			Str("error", err.Error()).
			Msg("Failed to auto-migrate database models")
		return nil, err
	}

	logger.LogInfoMessage().
		Str("host", cfg.DBHost).
		Str("port", cfg.DBPort).
		Str("database", cfg.DBName).
		Dur("connection_time", time.Since(start)).
		Msg("Database connected and migrated successfully")

	return db, nil
}
