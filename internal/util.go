package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func ResolvePath(path string) string {
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		return path
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to resolve path: %s: %s", path, err))
	}
	return filepath.Join(cwd, path)
}
