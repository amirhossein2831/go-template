package env

import (
	"os"
)

type ServiceMode string

const ServiceModeServer ServiceMode = "server"

func GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}
