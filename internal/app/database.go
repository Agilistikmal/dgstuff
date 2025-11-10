package app

import (
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(provider string, url string) *gorm.DB {
	var db *gorm.DB
	var err error

	switch provider {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(url), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	default:
		logrus.Errorf("invalid database provider: %s", provider)
	}

	if err != nil {
		logrus.Errorf("failed to connect to database: %v", err)
		return nil
	}

	err = db.AutoMigrate(&model.Stuff{})
	if err != nil {
		logrus.Errorf("failed to migrate database: %v", err)
		return nil
	}

	return db
}
