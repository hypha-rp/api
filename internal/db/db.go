package db

import (
	"fmt"
	"hypha/api/internal/config"
	"hypha/api/internal/db/tables"

	"github.com/go-orm/gorm"
	_ "github.com/go-orm/gorm/dialects/postgres"
	"github.com/rs/zerolog/log"
)

var (
	DBConn *gorm.DB
)

type DBConnWrapper struct {
	DB *gorm.DB
}

func Connect(cfg *config.Config) (*gorm.DB, error) {
	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("dbname", cfg.Database.Dbname).
		Msg("Attempting to connect to the database")

	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Dbname,
		cfg.Database.Password,
		cfg.Database.Sslmode,
	))
	if err != nil {
		log.Error().
			Err(err).
			Str("host", cfg.Database.Host).
			Int("port", cfg.Database.Port).
			Str("dbname", cfg.Database.Dbname).
			Msg("Failed to connect to the database")
		return nil, err
	}

	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("dbname", cfg.Database.Dbname).
		Msg("Successfully connected to the database")

	DBConn = db
	return DBConn, nil
}

var tables_slice = []interface{}{
	&tables.Product{},
	&tables.Integration{},
}

func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("Starting database migration")
	for _, table := range tables_slice {
		if err := db.AutoMigrate(table).Error; err != nil {
			log.Error().Err(err).Msg("Database migration failed")
			return err
		}
	}
	log.Info().Msg("Database migration completed successfully")
	return nil
}

func (wrapper *DBConnWrapper) Create(record interface{}) error {
	return wrapper.DB.Create(record).Error
}

func (wrapper *DBConnWrapper) First(result interface{}, conditions ...interface{}) error {
	return wrapper.DB.First(result, conditions...).Error
}
