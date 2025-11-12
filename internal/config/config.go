package config

import (
	_ "embed"
	"errors"
	"os"
	"strings"

	"github.com/agilistikmal/dgstuff/internal/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed config_template.yml
var configTemplate string

func LoadConfig() *viper.Viper {
	viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			logrus.Warnf("config file not found. creating default config")

			randomKey := pkg.GenerateRandomString(32)
			configTemplate = strings.Replace(configTemplate, "<random-key>", randomKey, 1)

			os.WriteFile("config.yml", []byte(configTemplate), 0644)
			logrus.Warnf("default config created. you can edit the config file and run the application again to apply the changes")

			if err := viper.ReadInConfig(); err != nil {
				logrus.Fatalf("failed to read config: %v", err)
			}

			return viper.GetViper()
		} else {
			logrus.Fatalf("failed to read config: %v", err)
		}
	}

	return viper.GetViper()
}
