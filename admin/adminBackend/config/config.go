package config

import (
	"log"

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
