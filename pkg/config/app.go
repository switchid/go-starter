package config

func (cfg *Config) GetAppName() string {
	if cfg != nil && cfg.App != nil {
		return cfg.App.Name
	}
	return ""
}

func (cfg *Config) GetAppVersion() string {
	if cfg != nil && cfg.App != nil {
		return cfg.App.Version
	}
	return ""
}

func (cfg *Config) GetAppDebug() bool {
	if cfg != nil && cfg.App != nil {
		return cfg.App.Debug
	}
	return false
}
