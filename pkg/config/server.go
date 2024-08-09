package config

import "strconv"

func (cfg *Config) GetServerHost() string {
	if cfg != nil && cfg.Server != nil {
		return cfg.Server.Host
	}
	return ""
}

func (cfg *Config) GetServerPort() string {
	if cfg != nil && cfg.Server != nil {
		return strconv.Itoa(cfg.Server.Port)
	}
	return ""
}
