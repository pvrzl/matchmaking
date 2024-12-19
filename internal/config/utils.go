package config

import (
	"os"
	"strconv"
)

func GetEnvBool(key string, dflt bool) bool {
	value := os.Getenv(key)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return dflt
	}
	return boolValue
}

func GetEnvString(key string, dflt string) string {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	return value
}

func GetEnvInt(key string, dflt int) int {
	value := os.Getenv(key)
	i, err := strconv.ParseInt(value, 10, 64)
	if value == "" && err != nil {
		return dflt
	}
	return int(i)
}

func GetEnvFloat(key string, dflt float64) float64 {
	value := os.Getenv(key)
	i, err := strconv.ParseFloat(value, 64)
	if value == "" && err != nil {
		return dflt
	}
	return i
}