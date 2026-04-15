package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port   int    `env:"PORT"`
	AppEnv string `env:"APP_ENV"`

	DatabaseURL string `env:"DB"`

	JWTSecret string `env:"JWT_SECRET"`

	AzureStorageURI string `env:"AZURE_STORAGE_URI"`
	AzureContainer  string `env:"AZURE_CONTAINER"`

	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
}

func Load() *Config {
	return &Config{
		Port:   getInt("PORT", 8080),
		AppEnv: getEnv("APP_ENV", "development"),

		DatabaseURL: getEnvRequired("DB"),

		JWTSecret: getEnv("JWT_SECRET", "secret-key"),

		AzureStorageURI: getEnvRequired("AZURE_STORAGE_URI"),
		AzureContainer:  getEnvRequired("AZURE_CONTAINER"),

		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ","),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	panic(fmt.Sprintf("Key (%s) not defined in .env!", key))
}

func getInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err != nil {
			return i
		}
	}
	return fallback
}
