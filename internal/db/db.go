package db

import (
	"fmt"
	"hypha/api/internal/config"

	"github.com/go-orm/gorm"
	_ "github.com/go-orm/gorm/dialects/postgres"
)

var (
	DBConn *gorm.DB
)

type DBConnWrapper struct {
	DB *gorm.DB
}

func Connect(cfg *config.Config) (*gorm.DB, error) {
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
		return nil, err
	}

	DBConn = db
	return DBConn, nil
}

func AutoMigrate(dbConn *gorm.DB, tables ...interface{}) error {
	for _, table := range tables {
		if err := dbConn.AutoMigrate(table).Error; err != nil {
			return err
		}
	}
	return nil
}

func (wrapper *DBConnWrapper) Create(record interface{}) error {
	return wrapper.DB.Create(record).Error
}

func (wrapper *DBConnWrapper) First(result interface{}, conditions ...interface{}) error {
	return wrapper.DB.First(result, conditions...).Error
}
