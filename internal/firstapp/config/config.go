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
type Config struct {
	APP      APPConfig      `mapstructure:"app"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// APPConfig hold the app properties
type APPConfig struct {
	Name string `mapstructure:"name"`
}

// LoggerConfig	hold the config for logger
type LoggerConfig struct {
	Level             string                `mapstructure:"level"`
	Encoding          string                `mapstructure:"encoding"`
	Development       string                `mapstructure:"development"`
	OutputPaths       string                `mapstructure:"outputPaths"`
	ErrorOutputPaths  string                `mapstructure:"errorOutputPaths"`
	DisableCaller     string                `mapstructure:"disableCaller"`
	DisableStacktrace string                `mapstructure:"disableStacktrace"`
	Sampling          *SamplingLoggerConfig `mapstructure:"sampling"`
}

// SamplingLoggerConfig	hold the config for sampling of logger
type SamplingLoggerConfig struct {
	Initial    string `mapstructure:"initial"`
	Thereafter string `mapstructure:"thereafter"`
}

// ServerConfig	hold the server config for http and grpc
type ServerConfig struct {
	HTTP HTTPServerConfig `mapstructure:"http"`
	GRPC GRPCServerConfig `mapstructure:"grpc"`
}

// HTTPServerConfig hold the http server config
type HTTPServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// GRPCServerConfig hold the grpc server config
type GRPCServerConfig struct {
	Port       string `mapstructure:"port"`
	Reflection string `mapstructure:"reflection"`
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

	// configure the viper
	viper.SetConfigName(configDir)
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

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

func (config *Config) GetEnv(key string) string {
	return env.GetExpandedEnv(key)
}
