package cluster

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFilePath(fileName string) string {
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	return filepath.Join(cwd, fileName)
}
