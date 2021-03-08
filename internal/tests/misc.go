package tests

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ReadFileToString - read the contents of a file and optionally trim any whitespace.
func ReadFileToString(filename string, trim bool) (contents string, err error) {
	var bytes []byte

	if bytes, err = os.ReadFile(filename); err != nil {
		return "", errors.Wrapf(err, "Failed to read file '%s'", filename)
	}

	contents = string(bytes)

	if trim {
		contents = strings.TrimSpace(contents)
	}

	return
}

// CreateTempFile - create a temp file with the filename pattern
func CreateTempFile(fileNamePattern string) (file *os.File, err error) {
	if file, err = os.CreateTemp("", fileNamePattern); err != nil {
		return nil, errors.Wrapf(err, "Failed to create temp file with pattern '%s'", fileNamePattern)
	}

	return
}

// CreateTempFileWithContents - create a temp file with the filename pattern and contents specified
func CreateTempFileWithContents(fileNamePattern, contents string) (file *os.File, err error) {
	if file, err = CreateTempFile(fileNamePattern); err != nil {
		return nil, err
	}

	if err = os.WriteFile(file.Name(), []byte(contents), 0666); err != nil {
		return nil, errors.Wrapf(err, "Failed to write contents to temp file '%s'", file.Name())
	}

	return
}
