package configs

import (
	"fmt"
	"log"
	"strings"

	"event-collector/pkg/env"

	"github.com/spf13/viper"
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

// APPConfig holds the app properties
type APPConfig struct {
	Name        string `mapstructure:"name"`
	ServiceMOde string `mapstructure:"service_mode"`
}

// LoggerConfig holds the config for logger
type LoggerConfig struct {
	Level             string                `mapstructure:"level"`
	Encoding          string                `mapstructure:"encoding"`
	Development       bool                  `mapstructure:"development"`
	OutputPaths       []string              `mapstructure:"outputPaths"`
	ErrorOutputPaths  []string              `mapstructure:"errorOutputPaths"`
	DisableCaller     bool                  `mapstructure:"disableCaller"`
	DisableStacktrace bool                  `mapstructure:"disableStacktrace"`
	Sampling          *SamplingLoggerConfig `mapstructure:"sampling"`
}

// SamplingLoggerConfig holds the config for sampling of logger
type SamplingLoggerConfig struct {
	Initial    int `mapstructure:"initial"`
	Thereafter int `mapstructure:"thereafter"`
}

// ServerConfig holds the server config for http and grpc
type ServerConfig struct {
	HTTP HTTPServerConfig `mapstructure:"http"`
	GRPC GRPCServerConfig `mapstructure:"grpc"`
}

// HTTPServerConfig holds the http server config
type HTTPServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// GRPCServerConfig holds the grpc server config
type GRPCServerConfig struct {
	Port       int  `mapstructure:"port"`
	Reflection bool `mapstructure:"reflection"`
}

// DatabaseConfig holds the database config
type DatabaseConfig struct {
	URI     string               `mapstructure:"uri"`
	Migrate MongoMigrationConfig `mapstructure:"migrate"`
}

// MongoMigrationConfig holds the database migration config
type MongoMigrationConfig struct {
	Path string `mapstructure:"path"`
	Type string `mapstructure:"type"`
	Step int    `mapstructure:"step"`
}

// NewConfig loads configuration from file and environment.
func NewConfig() (*Config, error) {
	appEnv := env.GetEnv("APP_ENV", "local")
	configName, ok := profileConfigFileMap[appEnv]
	if !ok {
		log.Printf("WARN: unknown APP_ENV=%q, falling back to local", appEnv)
		configName = profileConfigFileMap["local"]
	}

	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yml")
	v.AddConfigPath("./configs")

	// Enable env overrides like DATABASE_URI, SERVER_GRPC_PORT, etc.
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults (use dot notation for nested keys)
	v.SetDefault("database.uri", "mongodb://localhost:27017")
	v.SetDefault("server.http.host", "0.0.0.0")
	v.SetDefault("server.http.port", 8080)
	v.SetDefault("server.grpc.port", 9090)
	v.SetDefault("server.grpc.reflection", false)

	if err := v.ReadInConfig(); err != nil {
		// If you want config file to be optional, handle viper.ConfigFileNotFoundError here.
		return nil, fmt.Errorf("failed to read config file %q: %w", configName, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	log.Printf("Configuration loaded successfully! (env=%s, file=%s)", appEnv, configName)
	return &cfg, nil
}
