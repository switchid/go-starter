package config

import (
	"strconv"
	"strings"
)

func (cfg *Config) GetDBDriver() string {
	if cfg != nil && cfg.Database != nil {
		return strings.ToLower(cfg.Database.Driver)
	}
	return ""
}

func (cfg *Config) GetDBHost() string {
	if cfg != nil && cfg.Database != nil {
		return cfg.Database.Host
	}
	return ""
}

func (cfg *Config) GetDBPort() string {
	if cfg != nil && cfg.Database != nil {
		return strconv.Itoa(cfg.Database.Port)
	}
	return ""
}

func (cfg *Config) GetDBName() string {
	if cfg != nil && cfg.Database != nil {
		return cfg.Database.Name
	}
	return ""
}

func (cfg *Config) GetDBUser() string {
	if cfg != nil && cfg.Database != nil && cfg.Database.Credentials != nil {
		return cfg.Database.Credentials.Username
	}
	return ""
}

func (cfg *Config) GetDBPassword() string {
	if cfg != nil && cfg.Database != nil && cfg.Database.Credentials != nil {
		return cfg.Database.Credentials.Password
	}
	return ""
}

func (cfg *Config) GetDBPrefix() string {
	if cfg != nil && cfg.Database != nil {
		return cfg.Database.Prefix
	}
	return ""
}
