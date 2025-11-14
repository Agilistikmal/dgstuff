package app

import (
	"os/exec"

	"github.com/agilistikmal/dgstuff/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
		&model.Stuff{}, &model.StuffCategory{}, &model.StuffMedia{},
		&model.Stock{},
		&model.Transaction{}, &model.TransactionStuff{}, &model.TransactionPayment{},
		&model.AppInfo{},
	)
	if err != nil {
		logrus.Errorf("failed to migrate database: %v", err)
		return nil
	}

	var appInfo model.AppInfo
	db.First(&appInfo)

	// get version from git tag
	version, _ := exec.Command("git", "describe", "--tags").Output()
	if string(version) == "" {
		version = []byte("Latest")
	}

	err = db.Save(&model.AppInfo{
		ID:          appInfo.ID,
		Name:        viper.GetString("app.name"),
		Description: viper.GetString("app.description"),
		LogoURL:     viper.GetString("app.logo_url"),
		Version:     string(version),
		FirstLaunch: true,
	}).Error
	if err != nil {
		logrus.Errorf("failed to save app info: %v", err)
		return nil
	}

	return db
}
