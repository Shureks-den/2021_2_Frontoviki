package config

import "os"

type DatabaseConfig struct {
	DatabaseUrl string
}

type Config struct {
	DbConfig DatabaseConfig
}

func NewConfig() *Config {
	return &Config{
		DbConfig: DatabaseConfig{DatabaseUrl: GetEnv("DATABASE_URL", "")},
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
