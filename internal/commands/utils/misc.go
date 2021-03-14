package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// getEnv - utility function available in templates as "getenv"
func getEnv(key string) string {
	value, found := os.LookupEnv(key)

	if found {
		return value
	}
	return ""
}

// getFirstMatchedFile - from the given pattern, it turns the filename (without dir) of the first matching file
func getFirstMatchedFile(pattern string) (string, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(filenames) == 0 {
		return "", fmt.Errorf("No files matched for pattern: %s", pattern)
	}

	_, file := path.Split(filenames[0])
	return file, nil
}
