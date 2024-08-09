package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Flags string

type Flag interface {
	SetFlag(flag string, value string) string
}

func (f *Flags) SetFlag(flag, value string) string {
	*f = Flags(strings.TrimSpace(string(*f) + " " + fmt.Sprintf("-%s=%s", flag, value)))
	return string(*f)
}

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

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func FileExistsWithContains(dir string, contain string) bool {
	foundFile := false
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), contain) {
			foundFile = true
			return filepath.SkipDir // Exit early
		}
		return nil
	})

	if err != nil {
		return false
	}

	return foundFile
}
