package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetProjectPath() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Can't get current working directory\n")
		return "", err
	}
	parentDir := filepath.Dir(workDir)

	return strings.TrimSpace(parentDir), nil
}

func GetCurrentDirectory() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Base(workDir), nil
}

func GetCurrentExecutableDir() (string, error) {
	exeDir, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exeDir), nil
}

func GetCurrentExecutableName() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeName := filepath.Base(exePath)
	return exeName, nil
}
