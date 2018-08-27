package cluster

import (
	"os"
	"path/filepath"
)

// GetFilePath fetches and returns the current working directory.
func GetFilePath(fileName string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, fileName)
}
