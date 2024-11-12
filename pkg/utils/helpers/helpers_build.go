package helpers

import (
	"runtime/debug"
	"syscall"
)

func EnableANSI() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	var mode uint32
	consoleHandle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	err := syscall.GetConsoleMode(consoleHandle, &mode)
	if err != nil {
		return
	}

	// Enable ANSI escape sequence processing
	mode |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_, _, err = setConsoleMode.Call(uintptr(consoleHandle), uintptr(mode))
	if err != nil {
		return
	}
}

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
