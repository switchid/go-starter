package config

import (
	"GoStarter/pkg/utils/paths"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"path"
)

type Config struct {
	App *struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
		Key     string `mapstructure:"key,omitempty"`
		Debug   bool   `mapstructure:"debug,omitempty"`
	}
	Server *struct {
		App *struct {
			Host string
			Port int
		}
		Api *struct {
			Host string
			Port int
		}
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
	Services *struct {
		App *struct {
			Name        string
			Display     string
			Description string
		}
		Api *struct {
			Name        string
			Display     string
			Description string
		}
	} `mapstructure:",omitempty"`
}

func Load() (*Config, error) {
	appPath, errPath := paths.GetCurrentExecutableDir()
	if errPath != nil {
		return nil, fmt.Errorf("error unable get config path: %w", errPath)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(appPath + "/configs")

	// Set default values
	viper.SetDefault("app.debug", false)
	viper.SetDefault("app.key", "")
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.credentials.name", "root")
	viper.SetDefault("database.credentials.password", "")
	viper.SetDefault("logging.level", "warn")
	viper.SetDefault("server.app.host", "localhost")
	viper.SetDefault("server.app.port", 8080)
	viper.SetDefault("server.api.host", "localhost")
	viper.SetDefault("server.api.port", 9090)
	viper.SetDefault("services.app.name", "MyApp")
	viper.SetDefault("services.app.display", "MyApp")
	viper.SetDefault("services.app.description", "Service Application of MyApp")
	viper.SetDefault("services.api.name", "MyAppApi")
	viper.SetDefault("services.api.display", "MyApp (API)")
	viper.SetDefault("services.api.description", "Service Application of MyApp (API)")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}

func (cfg *Config) SaveConfig() error {
	appPath, errPath := paths.GetCurrentExecutableDir()
	if errPath != nil {
		return fmt.Errorf("error unable get config path: %w", errPath)
	}

	configMap := make(map[string]interface{})
	errMapping := mapstructure.Decode(cfg, &configMap)
	if errMapping != nil {
		return fmt.Errorf("error decoding config to map: %w", errMapping)
	}

	errMerge := viper.MergeConfigMap(configMap)
	if errMerge != nil {
		return fmt.Errorf("error decoding config to merge: %w", errMerge)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path.Join(appPath, "configs"))

	errSave := viper.WriteConfig()
	if errSave != nil {
		return fmt.Errorf("error save config file: %w", errSave)
	}

	return nil
}
