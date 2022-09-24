package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Appname  string
	Server   Server
	LogLevel string
}

type Server struct {
	Port string
}

func New(configPath, configName string) (Config, error) {
	viperConfig, err := ReadConfig(configPath, configName)
	if err != nil {
		return Config{}, err
	}
	config := Config{}

	if err := viperConfig.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func ReadConfig(configPath, configName string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	err := v.ReadInConfig()

	return v, err
}
