package paths

import (
	"os"
	"path/filepath"
	"strings"
)

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
