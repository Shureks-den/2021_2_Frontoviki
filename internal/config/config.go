package config

import (
	"flag"
	"fmt"
	"os"
	"yula/internal/pkg/logging"

	"github.com/spf13/viper"
)

type config struct {
	Server struct {
		Main struct {
			Host   string
			Port   string
			Secure bool
		}
	}

	Databases struct {
		Postgres struct {
			Host     string
			Port     string
			User     string
			Password string
			DbName   string
		}

		Tarantool struct {
			Host     string
			Port     string
			User     string
			Password string
		}
	}
}

var (
	logger logging.Logger = logging.GetLogger()
	Cfg    config
)

func argparse() (string, string, string) {
	mode := flag.String("mode", "dev", "config mode")
	filePath := flag.String("path", ".", "path to config file")
	configName := fmt.Sprintf("config_%s", *mode)
	return *mode, *filePath, configName
}

func LoadConfig() error {
	mode, configPath, configName := argparse()

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	logger.Infof("configuration file with mode=%s path=%s name=%s", mode, configPath, configName)

	if err := setConfig(); err != nil {
		return err
	}
	logger.Info(Cfg)
	return nil
}

func setConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&Cfg)
	return err
}

func (c *config) GetMainHost() string {
	return c.Server.Main.Host
}

func (c *config) GetMainPort() string {
	return c.Server.Main.Port
}

func (c *config) GetMainSchema() string {
	if c.Server.Main.Secure {
		return "https"
	}
	return "http"
}

func (c *config) IsSecure() bool {
	return Cfg.Server.Main.Secure
}

func (c *config) GetPostgresUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		c.Databases.Postgres.User,
		c.Databases.Postgres.Password,
		c.Databases.Postgres.Host,
		c.Databases.Postgres.Port,
		c.Databases.Postgres.DbName,
	)
}

func (c *config) GetTarantoolCfg() *TarantoolConfig {
	return &TarantoolConfig{
		TarantoolServerAddress: fmt.Sprintf("%s:%s", c.Databases.Tarantool.Host, c.Databases.Tarantool.Port),
		TarantoolOpts: TarantoolOptions{
			User: c.Databases.Tarantool.User,
			Pass: c.Databases.Tarantool.Password,
		},
	}
}

type TarantoolOptions struct {
	User string
	Pass string
}

type TarantoolConfig struct {
	TarantoolServerAddress string
	TarantoolOpts          TarantoolOptions
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
