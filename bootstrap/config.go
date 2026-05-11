package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

func SetupConfig(path string) error {
	return godotenv.Load(path)
}

func GetEnv(key string, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		return defaultValue
	}
	return result
}
