package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetEnvVar[T any](varName string, fallback ...T) (T, error) {
	valueStr := os.Getenv(varName)

	if valueStr == "" {
		if len(fallback) > 0 {
			return fallback[0], nil
		}
		return zeroValue[T](), fmt.Errorf("environment variable '%s' must be set", varName)
	}

	var result T
	switch any(result).(type) {
	case int:
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return handleFallback(varName, fallback, fmt.Errorf("must be an integer: %v", err))
		}
		return any(value).(T), nil
	case float64:
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return handleFallback(varName, fallback, fmt.Errorf("must be a float: %v", err))
		}
		return any(value).(T), nil
	case bool:
		value, err := strconv.ParseBool(valueStr)
		if err != nil {
			return handleFallback(varName, fallback, fmt.Errorf("must be a boolean: %v", err))
		}
		return any(value).(T), nil
	case string:
		return any(valueStr).(T), nil
	case time.Duration:
		value, err := time.ParseDuration(valueStr)
		if err != nil {
			return handleFallback(varName, fallback, fmt.Errorf("must be a duration: %v", err))
		}
		return any(value).(T), nil
	default:
		return zeroValue[T](), fmt.Errorf("unsupported type for environment variable '%s'", varName)
	}
}

func handleFallback[T any](varName string, fallback []T, err error) (T, error) {
	if len(fallback) > 0 {
		return fallback[0], nil
	}
	return zeroValue[T](), fmt.Errorf("environment variable '%s' %v", varName, err)
}

func zeroValue[T any]() T {
	var zero T
	return zero
}
