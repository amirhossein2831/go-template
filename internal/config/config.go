package config

import (
	"event-collector/pkg/env"
	"github.com/spf13/viper"
	"log"
)

var profileConfigFileMap = map[string]string{
	"local":       "application",
	"development": "application-development",
	"production":  "application-production",
	"test":        "application-test",
}

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

// NewConfig loads configuration from file and environment, returning a singleton instance.
func NewConfig() (*Config, error) {
	// read app env
	appEnv := env.GetEnv("APP_ENV", "local")
	configDir := profileConfigFileMap[appEnv]
	log.Println(configDir)

	viper.SetConfigName(configDir)   // Name of config file (without extension)
	viper.SetConfigType("yml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./configs") // Path to look for the config file in
	viper.AutomaticEnv()             // Read in environment variables that match

	// ---> SET DEFAULTS HERE <---
	// Use dot notation for nested keys
	viper.SetDefault("database.uri", "mongodb://localhost:27017")

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("FATAL: Failed to read Config: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("FATAL: Failed to unmarshal Config: %v", err)
		return nil, err
	}

	log.Println("Configuration loaded successfully!")
	return &config, nil
}
