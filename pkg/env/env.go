package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

func GetTypeEnv[O string | int | bool | float64](key string, defaultValue O) O {
	envStrValue := os.Getenv(key)
	if envStrValue == "" {
		return defaultValue
	}

	var value any
	var err error
	switch any(*new(O)).(type) {
	case float64:
		value, err = strconv.ParseFloat(envStrValue, 64)
	case int:
		value, err = strconv.Atoi(envStrValue)
	case bool:
		value, err = strconv.ParseBool(envStrValue)
	default:
		value = envStrValue
	}
	if err != nil {
		fmt.Println("GetConfig: invalid value for key: " + key)
		return defaultValue
	}
	oValue, _ := value.(O)
	return oValue
}
