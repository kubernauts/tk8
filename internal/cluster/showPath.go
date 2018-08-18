package cluster

import (
	"os"
	"path/filepath"
)

func GetFilePath(fileName string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, fileName)
}
