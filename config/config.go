package config

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	AdminPort    string `mapstructure:"ADMIN_PORT"`
// 	CustomerPort string `mapstructure:"CUSTOMER_PORT"`
// 	TraderPort   string `mapstructure:"TRADER_PORT"`

// 	DBHost     string `mapstructure:"DB_HOST"`
// 	DBPort     string `mapstructure:"DB_PORT"`
// 	DBUser     string `mapstructure:"DB_USER"`
// 	DBPassword string `mapstructure:"DB_PASSWORD"`
// 	DBName     string `mapstructure:"DB_NAME"`

// 	JWTSecret    string `mapstructure:"JWT_SECRET"`
// 	CookieDomain string `mapstructure:"COOKIE_DOMAIN"`

// 	AdminEmail    string `mapstructure:"Admin_Email"`
// 	AdminPassword string `mapstructure:"Admin_Password"`

// 	PolygonApiKey string `mapstructure:"POLYGON_API_KEY"`
// }

// var AppConfig Config

// func LoadConfig() (*Config, error) {
// 	wd, _ := os.Getwd()
// 	fmt.Printf("[DEBUG] Current Directory: %s\n", wd)
// 	fmt.Println("[DEBUG] Viper is looking for a '.env' file here.")

// 	fmt.Printf("[DEBUG] Current Working Directory: %s\n", wd)

// 	viper.SetConfigFile(".env")
// 	viper.AutomaticEnv() // Enable reading environment variables

// 	if err := viper.ReadInConfig(); err != nil {
// 		// Handle the case where .env file is not found gracefully,
// 		// especially if env vars are expected to take precedence
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			log.Println("WARNING: .env file not found, relying on environment variables.")
// 		} else {
// 			return nil, fmt.Errorf("failed to read config file: %w", err)
// 		}
// 	}

// 	var cfg Config
// 	if err := viper.Unmarshal(&cfg); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
// 	}

// 	fmt.Println("[DEBUG] Viper has successfully loaded configuration.")
// 	log.Println("Configuration loaded successfully.")
// 	AppConfig = cfg // Assign to global AppConfig
// 	return &cfg, nil
// }
