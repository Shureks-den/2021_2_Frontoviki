package config

import "os"

type DatabaseConfig struct {
	DatabaseUrl string
}

type TarantoolOptions struct {
	User string
	Pass string
}

type TarantoolConfig struct {
	TarantoolServerAddress string
	TarantoolOpts          TarantoolOptions
}

type Config struct {
	DbConfig     DatabaseConfig
	TarantoolCfg TarantoolConfig
}

func NewConfig() *Config {
	return &Config{
		DbConfig: DatabaseConfig{DatabaseUrl: GetEnv("DATABASE_URL", "")},
		TarantoolCfg: TarantoolConfig{
			TarantoolServerAddress: GetEnv("TARANTOOL_ADDRESS", "localhost:3302"),
			TarantoolOpts: TarantoolOptions{
				User: GetEnv("TARANTOOL_USER", "admin"),
				Pass: GetEnv("TARANTOOL_PASS", "pass"),
			},
		},
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
