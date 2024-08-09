package loggers

import (
	"GoStarter/pkg/utils/helpers"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	*log.Logger
}

type LoggerPrint interface {
	LogInfo(format string, v ...interface{})
	LogError(format string, v ...interface{})
}

func NewLogger() (*Logger, error) {
	exePath, errPath := helpers.GetCurrentExecutableDir()
	if errPath != nil {
		log.Fatalf("Error getting executable path: %v", errPath)
	}

	logPath := filepath.Join(exePath, "log")
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logPath, fmt.Sprintf("service-%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	l := log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &Logger{l}, nil
}

func (s *Logger) LogError(format string, v ...interface{}) {
	if s != nil {
		s.SetPrefix("ERROR: ")
		s.Printf(format, v...)
	}
}

func (s *Logger) LogInfo(format string, v ...interface{}) {
	if s != nil {
		s.SetPrefix("INFO: ")
		s.Printf(format, v...)
	}
}
