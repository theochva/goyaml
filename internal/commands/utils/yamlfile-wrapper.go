package utils

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/theochva/goyaml/pkg/yamldoc"
	"github.com/theochva/goyaml/pkg/yamlfile"
)

// YamlFileWrapper - simple wrapper that can either read from STDIN and
// write to STDOUT when "saved" or normally to/from a file
type YamlFileWrapper struct {
	yamlfile.YamlFile
	pipeMode bool
	stdin    io.Reader
	stdout   io.Writer
}

// NewYamlFileWrapper - create new YamlFile wrapper
func NewYamlFileWrapper(filename string, stdin io.Reader, stdout io.Writer) yamlfile.YamlFile {
	var wrapperFilename = filename

	if filename == "" {
		wrapperFilename = "__PIPE__"
	}
	return &YamlFileWrapper{
		pipeMode: (filename == ""),
		YamlFile: yamlfile.New(wrapperFilename),
		stdin:    stdin,
		stdout:   stdout,
	}
}

// Exists - override
func (y *YamlFileWrapper) Exists() bool {
	if y.pipeMode {
		return true
	}
	return y.YamlFile.Exists()
}

// Filename - override
func (y *YamlFileWrapper) Filename() string {
	if y.pipeMode {
		return ""
	}
	return y.YamlFile.Filename()
}

// Load - override
func (y *YamlFileWrapper) Load() (loaded bool, err error) {
	if y.pipeMode {
		var yaml yamldoc.YamlDoc

		if yaml, err = yamldoc.New(y.stdin); err != nil {
			err = errors.Wrap(err, "Failed to read/parse yaml from stdin")
			return
		}
		y.YamlFile.SetData(yaml.Data())
	}
	return y.YamlFile.Load()
}

// LoadReader - override
func (y *YamlFileWrapper) LoadReader(reader io.Reader) (loaded bool, err error) {
	if y.pipeMode {
		var yaml yamldoc.YamlDoc

		if yaml, err = yamldoc.New(reader); err != nil {
			err = errors.Wrap(err, "Failed to read/parse yaml from stdin")
			return
		}
		y.YamlFile.SetData(yaml.Data())
	}
	return y.YamlFile.Load()
}

// Save - override
func (y *YamlFileWrapper) Save() (err error) {
	if y.pipeMode {
		var text string
		if text, err = y.Text(); err != nil {
			return errors.Wrap(err, "Failed to generate yaml text")
		}
		fmt.Fprintln(y.stdout, text)
		return
	}
	return y.YamlFile.Save()
}
