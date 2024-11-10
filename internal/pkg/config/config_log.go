package config

func (cfg *Config) GetLogLevel() string {
	if cfg != nil && cfg.Logging != nil {
		return cfg.Logging.Level
	}
	return ""
}

func (cfg *Config) SetLogLevel(level string) {
	if cfg != nil && cfg.Logging != nil {
		cfg.Logging.Level = level
	}
}

func (cfg *Config) GetLogFile() string {
	if cfg != nil && cfg.Logging != nil {
		return cfg.Logging.File
	}
	return ""
}

func (cfg *Config) SetLogFile(file string) {
	if cfg != nil && cfg.Logging != nil {
		cfg.Logging.File = file
	}
}
