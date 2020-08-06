package config

import (
	"os"
	"strconv"
	"strings"
)

// EnvGetInt Get Integer Config data
func EnvGetInt(key string, defaultVal int) int {
	data := os.Getenv(key)
	if len(data) == 0 {
		return defaultVal
	}
	out, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		return defaultVal
	}
	return int(out)
}

// EnvGetStr Get String Config data
func EnvGetStr(key string, defaultVal string) string {
	data := os.Getenv(key)
	if len(data) == 0 {
		return defaultVal
	}
	return data
}

// EnvGetBool Get Bool Config data
func EnvGetBool(key string, defaultVal bool) bool {
	data := os.Getenv(key)
	if len(data) == 0 {
		return defaultVal
	}
	if strings.ToLower(data) == "true" {
		return true
	}
	return false
}
