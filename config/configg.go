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
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
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

	v.SetConfigFile("config/config.yaml")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	log.Println("Config loaded from:", v.ConfigFileUsed())

	if err := v.Unmarshal(&AppConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&AppConfig); err != nil {
		return nil, err
	}

	return &AppConfig, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("server.admin_port", "8080")
	v.SetDefault("jwt.expire_hours", 24)
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if cfg.Database.Port == "" {
		return fmt.Errorf("database.port is required")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database.user is required")
	}
	if cfg.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}
	return nil
}
