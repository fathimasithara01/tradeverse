package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config struct matches the keys in your .env file
type Config struct {
	Port         string `mapstructure:"PORT"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	CookieDomain string `mapstructure:"COOKIE_DOMAIN"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read .env file: %v", err)
	}

	// Unmarshal the config into the global AppConfig variable
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	log.Println("Configuration loaded successfully.")
}
