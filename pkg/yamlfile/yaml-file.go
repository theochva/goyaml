package yamlfile

import (
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/theochva/goyaml/pkg/yamldoc"
)

// YamlFile - represents a yaml file
type YamlFile interface {
	yamldoc.YamlDoc

	// Exists - Check whether the file actually exists
	Exists() bool
	// Filename - returns the filename
	Filename() string
	// Load - loads the file (if it exists)
	Load() (loaded bool, err error)
	// LoadReader - load from a reader
	LoadReader(reader io.Reader) (loaded bool, err error)
	// Save - saves the yaml file
	Save() (err error)
}

type yamlFile struct {
	yamldoc.YamlDoc
	filename string
}

// New - create a new yaml YamlFile object. Returns nil if filename is empty
func New(filename string) YamlFile {
	if filename == "" {
		return nil
	}

	result := &yamlFile{
		filename: filename,
	}

	result.YamlDoc, _ = yamldoc.New(nil)

	return result
}

// Load - create a new YamlFile and load YAML contents.
//
// If the file does not exist, then the YamlFile is empty and loaded=false.
//
// If the file exists and parsed successfully, then YamlFile is populated with the contents
// of the YAML file and loaded=true.
//
// If an error occurs while opening or parsing the file then YamlFile=nil and "err" will
// contain the error information.
func Load(filename string) (loaded bool, yamlFile YamlFile, err error) {
	yamlFile = New(filename)

	loaded, err = yamlFile.Load()
	return
}

// Exists - Check whether the file actually exists
func (y *yamlFile) Exists() bool {
	return fileExists(y.filename)
}

// Filename - returns the filename
func (y *yamlFile) Filename() string {
	return y.filename
}

// Load - loads the file (if it exists)
//
// If the file does not exist, then it returns loaded=false.
//
// If the file exists and parsed successfully, then YamlFile is populated with the contents
// of the YAML file and loaded=true.
//
// If an error occurs while opening or parsing the file then YamlFile is unchanged and "err" will
// contain the error information.
func (y *yamlFile) Load() (loaded bool, err error) {
	// If the file does not exist at this point,
	if !y.Exists() {
		return
	}

	// Otherwise, we will attempt to load the contents of the file
	var (
		file *os.File
	)
	// Open file
	if file, err = os.Open(y.filename); err != nil {
		return
	}

	defer file.Close()

	return y.LoadReader(file)
}

// LoadReader - load from a reader
func (y *yamlFile) LoadReader(reader io.Reader) (loaded bool, err error) {
	if reader != nil {
		if y.YamlDoc, err = yamldoc.New(reader); err != nil {
			return false, err
		}

		return true, nil
	}
	return false, nil
}

// Save - saves the yaml file
func (y *yamlFile) Save() (err error) {
	yamlBytes, err := y.Bytes()
	if err != nil {
		return errors.Wrap(err, "Failed to get yaml bytes")
	}
	return os.WriteFile(y.filename, yamlBytes, 0644)
}
