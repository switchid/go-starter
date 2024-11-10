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

func (cfg *Config) SetDBDriver(driver string) {
	if cfg != nil && cfg.Database != nil {
		cfg.Database.Driver = strings.ToLower(driver)
	}
}

func (cfg *Config) GetDBHost() string {
	if cfg != nil && cfg.Database != nil {
		return cfg.Database.Host
	}
	return ""
}

func (cfg *Config) SetDBHost(host string) {
	if cfg != nil && cfg.Database != nil {
		cfg.Database.Host = strings.ToLower(host)
	}
}

func (cfg *Config) GetDBPort() string {
	if cfg != nil && cfg.Database != nil {
		return strconv.Itoa(cfg.Database.Port)
	}
	return ""
}

func (cfg *Config) SetDBPort(port int) {
	if cfg != nil && cfg.Database != nil {
		cfg.Database.Port = port
	}
}

func (cfg *Config) GetDBName() string {
	if cfg != nil && cfg.Database != nil {
		return cfg.Database.Name
	}
	return ""
}

func (cfg *Config) SetDBName(name string) {
	if cfg != nil && cfg.Database != nil {
		cfg.Database.Name = name
	}
}

func (cfg *Config) GetDBUser() string {
	if cfg != nil && cfg.Database != nil && cfg.Database.Credentials != nil {
		return cfg.Database.Credentials.Username
	}
	return ""
}

func (cfg *Config) SetDBUser(user string) {
	if cfg != nil && cfg.Database != nil && cfg.Database.Credentials != nil {
		cfg.Database.Credentials.Username = user
	}
}

func (cfg *Config) GetDBPassword() string {
	if cfg != nil && cfg.Database != nil && cfg.Database.Credentials != nil {
		return cfg.Database.Credentials.Password
	}
	return ""
}

func (cfg *Config) SetDBPassword(password string) {
	if cfg != nil && cfg.Database != nil && cfg.Database.Credentials != nil {
		cfg.Database.Credentials.Password = password
	}
}

func (cfg *Config) GetDBPrefix() string {
	if cfg != nil && cfg.Database != nil {
		return cfg.Database.Prefix
	}
	return ""
}

func (cfg *Config) SetDBPrefix(prefix string) {
	if cfg != nil && cfg.Database != nil {
		cfg.Database.Prefix = prefix
	}
}
