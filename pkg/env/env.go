package env

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`\${(.*?)}`)

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

// GetExpandedEnv finds placeholders like ${VAR:default} in a string and replaces them.
func GetExpandedEnv(value string) string {
	return placeholderRegex.ReplaceAllStringFunc(value, func(match string) string {
		parts := placeholderRegex.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		content := parts[1]
		var envVar, defaultValue string

		if colonIndex := strings.Index(content, ":"); colonIndex != -1 {
			envVar = content[:colonIndex]
			defaultValue = content[colonIndex+1:]
		} else {
			envVar = content
		}

		if v, exists := os.LookupEnv(envVar); exists {
			return v
		}
		return defaultValue
	})
}
