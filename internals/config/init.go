package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func GetTestConfigFromFile() *Config {
	config, err := fetchConfigFromFile("test")
	if err != nil {
		fmt.Printf("couldn't read config: %s", err)
		os.Exit(1)
	}
	return config
}

func GetConfigFromFile() *Config {
	config, err := fetchConfigFromFile("")
	if err != nil {
		fmt.Printf("couldn't read config: %s", err)
		os.Exit(1)
	}
	return config
}

func fetchConfigFromFile(fileName string) (*Config, error) {
	if fileName == "" {
		fileName = "default"
	}

	viper.SetConfigName(fileName)
	viper.AddConfigPath("../conf/")
	viper.AddConfigPath("../../conf/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf/")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func WatchConfigChanges(l *zerolog.Logger, c *Config) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		newConfig, err := fetchConfigFromFile("")
		if err != nil {
			l.Err(err).Interface("current_config", c).Msg("watching config changes failed")
		} else {
			c = newConfig
			l.Debug().Interface("updated_config", c).Send()
		}
	})
	viper.WatchConfig()
}
