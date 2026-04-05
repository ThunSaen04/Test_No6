package database

import (
	"fmt"
	"log"

	"no6/backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewPostgres(cfg config.Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Warn),
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN()), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	log.Println("database connected")
	return db, nil
}
