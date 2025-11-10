package main

import (
	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/agilistikmal/dgstuff/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.LoadConfig()

	db := app.NewDatabase(viper.GetString("database.provider"), viper.GetString("database.url"))
	_ = app.NewValidator()

	logrus.Infof("Database %s and validator initialized", db.Name())
}
