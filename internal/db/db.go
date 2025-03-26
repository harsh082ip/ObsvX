package db

import (
	"fmt"
	"log"

	"github.com/harsh082ip/ObsvX/config"
	"github.com/harsh082ip/ObsvX/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
		return nil, err
	}

	// Auto Migrate tables
	db.AutoMigrate(&models.Order{})
	return db, nil
}
