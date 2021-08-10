package config

import (
	"os"
	"strconv"
)
func getEnvStr(key string, oldVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return oldVal
	}
	return val
}

func getEnvInt(key string, oldVal int) int{
	data := os.Getenv(key)
	if data == "" {
		return oldVal
	}
	val, _ := strconv.Atoi(data)
	return val
}
func getEnvBool(key string, oldVal bool) bool {
	data := os.Getenv(key)
	if data == "" {
		return oldVal
	}
	val, _ := strconv.ParseBool(data)
	return val
}
