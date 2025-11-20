package app

import (
	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(provider string, url string) *gorm.DB {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default,
	}

	switch provider {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(url), gormConfig)
	case "mysql":
		db, err = gorm.Open(mysql.Open(url), gormConfig)
	case "postgres":
		db, err = gorm.Open(postgres.Open(url), gormConfig)
	default:
		logrus.Errorf("invalid database provider: %s", provider)
	}

	if err != nil {
		logrus.Errorf("failed to connect to database: %v", err)
		return nil
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Stuff{}, &model.StuffCategory{}, &model.StuffMedia{},
		&model.Stock{},
		&model.Transaction{}, &model.TransactionStuff{}, &model.TransactionStuffData{}, &model.TransactionPayment{},
	)
	if err != nil {
		logrus.Errorf("failed to migrate database: %v", err)
		return nil
	}

	return db
}
