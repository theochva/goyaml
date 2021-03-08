package commands

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/theochva/goyaml/yamldoc"
	"github.com/theochva/goyaml/yamlfile"
)

type _YamlFileWrapper struct {
	yamlfile.YamlFile
	pipeMode bool
}

func newYamlFileWrapper(filename string) yamlfile.YamlFile {
	if filename == "" {
		return &_YamlFileWrapper{
			pipeMode: true,
			YamlFile: yamlfile.New("__PIPE__"),
		}
	}
	return &_YamlFileWrapper{
		pipeMode: false,
		YamlFile: yamlfile.New(filename),
	}
}

func (y *_YamlFileWrapper) Exists() bool {
	if y.pipeMode {
		return true
	}
	return y.YamlFile.Exists()
}

func (y *_YamlFileWrapper) Filename() string {
	if y.pipeMode {
		return ""
	}
	return y.YamlFile.Filename()
}

func (y *_YamlFileWrapper) Load() (loaded bool, err error) {
	if y.pipeMode {
		var yaml yamldoc.YamlDoc

		if yaml, err = yamldoc.New(os.Stdin); err != nil {
			err = errors.Wrap(err, "Failed to read/parse yaml from stdin")
			return
		}
		y.YamlFile.SetData(yaml.Data())
	}
	return y.YamlFile.Load()
}

func (y *_YamlFileWrapper) Save() (err error) {
	if y.pipeMode {
		var text string
		if text, err = y.Text(); err != nil {
			return errors.Wrap(err, "Failed to generate yaml text")
		}
		fmt.Println(text)
		return
	}
	return y.YamlFile.Save()
}
