package config

func (cfg *Config) GetLogLevel() string {
	if cfg != nil && cfg.Logging != nil {
		return cfg.Logging.Level
	}
	return ""
}

func (cfg *Config) GetLogFile() string {
	if cfg != nil && cfg.Logging != nil {
		return cfg.Logging.File
	}
	return ""
}
