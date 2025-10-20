package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name   string
		Env    string
		Secret string
	}

	Admin struct {
		UserID     int
		Email      string
		Password   string
		Commission int
	}

	Server struct {
		AdminPort    string `mapstructure:"admin_port"`
		CustomerPort string `mapstructure:"customer_port"`
		TraderPort   string `mapstructure:"trader_port"`
	}

	Cookie struct {
		Domain string
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	JWT struct {
		Secret      string
		ExpireHours int `mapstructure:"expire_hours"`
	}

	Redis struct {
		Host string
		Port int
	}
}

var AppConfig Config

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := v.Unmarshal(&AppConfig); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	log.Printf("Loaded configuration for environment: %s", AppConfig.App.Env)
	return &AppConfig, nil
}
