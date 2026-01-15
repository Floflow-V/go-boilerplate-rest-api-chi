package database

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"go-boilerplate-rest-api-chi/internal/config"
	"go-boilerplate-rest-api-chi/internal/entity"
)

type Database struct {
	Gorm  *gorm.DB
	sqlDB *sql.DB
}

func Init(cfg config.Config, logger zerolog.Logger) (*Database, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         gormLogger.Default.LogMode(gormLogger.Silent),
		TranslateError: true,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to the database")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to find database instance")
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)

	if err := db.AutoMigrate(
		// Models
		&entity.Book{},
		&entity.Author{},
	); err != nil {
		logger.Error().Err(err).Msg("auto-migration failed")
		return nil, err
	}

	return &Database{
		Gorm:  db,
		sqlDB: sqlDB,
	}, nil
}

func (d *Database) Close() error {
	if d.sqlDB != nil {
		return d.sqlDB.Close()
	}
	return nil
}
