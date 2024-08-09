package config

import (
	"GoStarter/pkg/utils/helpers"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App *struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
		Debug   bool   `mapstructure:"debug,omitempty"`
	}
	Server *struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	Database *struct {
		Driver      string
		Host        string
		Port        int
		Name        string
		Prefix      string
		Credentials *struct {
			Username string
			Password string
		}
		Pool *struct {
			MaxConnections int `mapstructure:"max_connections"`
			IdleTimeout    int `mapstructure:"idle_timeout"`
		} `mapstructure:",omitempty"`
	}
	Features map[string]bool
	Logging  *struct {
		Level string
		File  string
	} `mapstructure:",omitempty"`
}

func Load() (*Config, error) {
	appPath, errPath := helpers.GetCurrentExecutableDir()
	if errPath != nil {
		return nil, fmt.Errorf("error unable get config path: %w", errPath)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(appPath + "/configs")

	// Set default values
	viper.SetDefault("app.debug", false)
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.credentials.name", "root")
	viper.SetDefault("database.credentials.password", "")
	viper.SetDefault("logging.level", "warn")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
