package config

func (cfg *Config) GetServiceAppName() string {
	if cfg != nil && cfg.Services != nil {
		return cfg.Services.App.Name
	}
	return ""
}

func (cfg *Config) SetServiceAppName(name string) {
	if cfg != nil && cfg.Services != nil {
		cfg.Services.App.Name = name
	}
}

func (cfg *Config) GetServiceAppDisplay() string {
	if cfg != nil && cfg.Services != nil {
		return cfg.Services.App.Display
	}
	return ""
}

func (cfg *Config) SetServiceAppDisplay(display string) {
	if cfg != nil && cfg.Services != nil {
		cfg.Services.App.Display = display
	}
}

func (cfg *Config) GetServiceAppDescription() string {
	if cfg != nil && cfg.Services != nil {
		return cfg.Services.App.Description
	}
	return ""
}

func (cfg *Config) SetServiceAppDescription(description string) {
	if cfg != nil && cfg.Services != nil {
		cfg.Services.App.Description = description
	}
}

func (cfg *Config) GetServiceApiName() string {
	if cfg != nil && cfg.Services != nil {
		return cfg.Services.Api.Name
	}
	return ""
}

func (cfg *Config) SetServiceApiName(name string) {
	if cfg != nil && cfg.Services != nil {
		cfg.Services.Api.Name = name
	}
}

func (cfg *Config) GetServiceApiDisplay() string {
	if cfg != nil && cfg.Services != nil {
		return cfg.Services.Api.Display
	}
	return ""
}

func (cfg *Config) SetServiceApiDisplay(display string) {
	if cfg != nil && cfg.Services != nil {
		cfg.Services.Api.Display = display
	}
}

func (cfg *Config) GetServiceApiDescription() string {
	if cfg != nil && cfg.Services != nil {
		return cfg.Services.Api.Description
	}
	return ""
}

func (cfg *Config) SetServiceApiDescription(description string) {
	if cfg != nil && cfg.Services != nil {
		cfg.Services.Api.Description = description
	}
}
