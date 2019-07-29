package main

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds the configuration
type Config struct {
	Port string
}

// LoadConfig is used to create and load the config
func LoadConfig() (*Config, error) {
	viper.SetConfigType("toml")
	viper.SetConfigFile("badge.conf")
	viper.AddConfigPath(".") // Search the root directory for the configuration file
	viper.SetDefault("Port", "COM1")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	var config Config
	err := viper.Unmarshal(&config)

	return &config, err
}
