package config

func (cfg *Config) GetAppName() string {
	if cfg != nil && cfg.App != nil {
		return cfg.App.Name
	}
	return ""
}

func (cfg *Config) SetAppName(name string) {
	if cfg != nil && cfg.App != nil {
		cfg.App.Name = name
	}
}

func (cfg *Config) GetAppVersion() string {
	if cfg != nil && cfg.App != nil {
		return cfg.App.Version
	}
	return ""
}

func (cfg *Config) SetAppVersion(version string) {
	if cfg != nil && cfg.App != nil {
		cfg.App.Version = version
	}
}

func (cfg *Config) GetAppDebug() bool {
	if cfg != nil && cfg.App != nil {
		return cfg.App.Debug
	}
	return false
}

func (cfg *Config) SetAppDebug(debug bool) {
	if cfg != nil && cfg.App != nil {
		cfg.App.Debug = debug
	}
}

func (cfg *Config) GetAppKey() string {
	if cfg != nil && cfg.App != nil {
		return cfg.App.Key
	}
	return ""
}

func (cfg *Config) SetAppKey(key string) {
	if cfg != nil && cfg.App != nil {
		cfg.App.Key = key
	}
}
