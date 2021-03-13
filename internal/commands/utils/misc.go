package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// FileOrDirectoryExists - check if specified filename (which is a file or directory) exists
func FileOrDirectoryExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDirectory - check if specified filename is a directory
func IsDirectory(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

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
