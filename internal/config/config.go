package config

import (
	_ "embed"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed config_template.yml
var configTemplate []byte

func LoadConfig() *viper.Viper {
	viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			logrus.Warnf("config file not found. creating default config")
			os.WriteFile("config.yml", configTemplate, 0644)
			logrus.Infof("default config created. please edit the config file and run the application again")
			os.Exit(0)
		} else {
			logrus.Fatalf("failed to read config: %v", err)
		}
	}

	return viper.GetViper()
}
