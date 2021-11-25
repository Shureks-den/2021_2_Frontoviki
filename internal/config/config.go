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

	Microservices struct {
		Chat struct {
			Host string
			Port string
		}

		Auth struct {
			Host string
			Port string
		}

		Category struct {
			Host string
			Port string
		}
	}

	Certificates struct {
		Selfsigned struct {
			Crt struct {
				Path string
			}

			Key struct {
				Path string
			}
		}

		Https struct {
			Crt struct {
				Path string
			}

			Key struct {
				Path string
			}
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

	flag.Parse()

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

func (c *config) GetChatEndPoint() string {
	return fmt.Sprintf("%s:%s", c.Microservices.Chat.Host, c.Microservices.Chat.Port)
}

func (c *config) GetAuthEndPoint() string {
	fmt.Println(c.Microservices.Auth.Host)
	fmt.Println(c.Microservices.Auth.Port)

	return fmt.Sprintf("%s:%s", c.Microservices.Auth.Host, c.Microservices.Auth.Port)
}

func (c *config) GetCategoryEndPoint() string {
	return fmt.Sprintf("%s:%s", c.Microservices.Category.Host, c.Microservices.Category.Port)
}

func (c *config) GetSelfSignedCrt() string {
	return c.Certificates.Selfsigned.Crt.Path
}

func (c *config) GetSelfSignedKey() string {
	return c.Certificates.Selfsigned.Key.Path
}

func (c *config) GetHTTPSCrt() string {
	return c.Certificates.Https.Crt.Path
}

func (c *config) GetHTTPSKey() string {
	return c.Certificates.Https.Key.Path
}
