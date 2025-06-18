package parse

import (
	"fmt"
	"strconv"
	"strings"
)

func ToPrimary[T string | int | int64 | float64 | bool](value string) (T, error) {
	var zero T

	switch any(zero).(type) {
	case string:
		return any(value).(T), nil
	case int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return zero, err
		}
		return any(i).(T), nil
	case int64:
		i64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return zero, err
		}
		return any(i64).(T), nil
	case float64:
		f64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return zero, err
		}
		return any(f64).(T), nil
	case bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return zero, err
		}
		return any(b).(T), nil
	default:
		return zero, fmt.Errorf("unsupported type for conversion: %T", zero)
	}
}

func ToStringArray(input string) []string {
	return strings.Split(input, ",")
}
