package storage

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"upisettle/internal/config"
	upilog "upisettle/internal/logger"
)

// NewDB initializes a GORM PostgreSQL connection.
func NewDB(cfg config.Config, log upilog.Logger) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Printf("database connection established")
	return db, nil
}

