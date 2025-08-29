package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port       string `mapstructure:"PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	JWTSecret    string `mapstructure:"JWT_SECRET"`
	CookieDomain string `mapstructure:"COOKIE_DOMAIN"`

	AdminEmail    string `mapstructure:"Admin_Email"`
	AdminPassword string `mapstructure:"Admin_Password"`

	PolygonApiKey string `mapstructure:"POLYGON_API_KEY"`
}

var AppConfig Config

func LoadConfig() (*Config, error) {
	wd, _ := os.Getwd()
	fmt.Printf("[DEBUG] Current Directory: %s\n", wd)
	fmt.Println("[DEBUG] Viper is looking for a '.env' file here.")

	fmt.Printf("[DEBUG] Current Working Directory: %s\n", wd)

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	fmt.Println("[DEBUG] Viper has successfully read the .env file.")
	fmt.Printf("[DEBUG] Value found for PORT is: '%s'\n", cfg.Port)

	log.Println("Configuration loaded successfully.")
	return &cfg, nil
}
