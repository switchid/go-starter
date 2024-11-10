package config

import "strconv"

func (cfg *Config) GetServerAppHost() string {
	if cfg != nil && cfg.Server != nil {
		if cfg.App != nil {
			return cfg.Server.App.Host
		}
		return ""
	}
	return ""
}

func (cfg *Config) SetServerAppHost(host string) {
	if cfg != nil && cfg.Server != nil {
		cfg.Server.App.Host = host
	}
}

func (cfg *Config) GetServerAppPort() string {
	if cfg != nil && cfg.Server != nil {
		if cfg.Server.App != nil {
			return strconv.Itoa(cfg.Server.App.Port)
		}
		return ""
	}
	return ""
}

func (cfg *Config) SetServerAppPort(port int) {
	if cfg != nil && cfg.Server != nil {
		cfg.Server.App.Port = port
	}
}

func (cfg *Config) GetServerApiHost() string {
	if cfg != nil && cfg.Server != nil {
		if cfg.Server.Api != nil {
			return cfg.Server.Api.Host
		}
		return ""
	}
	return ""
}

func (cfg *Config) SetServerApiHost(host string) {
	if cfg != nil && cfg.Server != nil {
		cfg.Server.Api.Host = host
	}
}

func (cfg *Config) GetServerApiPort() string {
	if cfg != nil && cfg.Server != nil {
		if cfg.Server.Api != nil {
			return strconv.Itoa(cfg.Server.Api.Port)
		}
		return ""
	}
	return ""
}

func (cfg *Config) SetServerApiPort(port int) {
	if cfg != nil && cfg.Server != nil {
		cfg.Server.Api.Port = port
	}
}
