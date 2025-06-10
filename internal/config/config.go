package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	once   sync.Once
	config *Config
)

// Config struct holds all configuration for the application.
// The tags `mapstructure:"..."` are used by Viper to unmarshal the data.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig	hold the server config
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig hold the database config
type DatabaseConfig struct {
	URI string `mapstructure:"uri"`
}

// GetConfig loads configuration from file and environment, returning a singleton instance.
// TODO: load base on active profile
func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("application") // Name of config file (without extension)
		viper.SetConfigType("yml")         // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("./configs")   // Path to look for the config file in
		viper.AutomaticEnv()               // Read in environment variables that match

		// ---> SET DEFAULTS HERE <---
		// Use dot notation for nested keys
		viper.SetDefault("database.uri", "mongodb://localhost:27017")

		// Read config
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}

		log.Println("Configuration loaded successfully!")
	})

	return config
}
