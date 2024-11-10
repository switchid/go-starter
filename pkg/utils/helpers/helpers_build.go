package helpers

import (
	"runtime/debug"
)

func DevMode() bool {
	// Check build info for custom tags
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			//log.Print(setting.Value + " " + setting.Key)
			if setting.Key == "-tags" && setting.Value == "dev" {
				return true
			}
		}
	}
	return false
}
