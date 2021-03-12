package commands

import (
	"github.com/pkg/errors"
	"github.com/theochva/goyaml/commands/cli"
	"github.com/theochva/goyaml/yamlfile"
)

// ValidationErrorAwareCommand - sub commands implement this interface when they
// want to be aware of YAML validation errors
type ValidationErrorAwareCommand interface {
	// IsValidationAware - return true if it needs to be validation aware
	IsValidationAware() bool
}

// SkipParsingCommand - sub commands implement this interface when they
// want to skip the initial parsing of the YAML file (i.e. do their own parsing)
type SkipParsingCommand interface {
	// ShouldSkipParsing - whether the input YAML should not be parsed up-front but parsed by the subcommand
	ShouldSkipParsing() bool
}

// GlobalOptions - global options that are needed to be accessed by all subcommands
type GlobalOptions interface {
	// YamlFile - get the yaml file we are working with
	YamlFile() yamlfile.YamlFile
	// IsPipe - whether we are reading the YAML from stdin
	IsPipe() bool
	// ValidationError - get the YAML validation error (if any)
	ValidationError() error
	// Load - loads the YAML file and any error is returned and also set in "ValidationError"
	Load() error
}

type _GlobalOptions struct {
	pipe              bool
	yamlFile          yamlfile.YamlFile
	yamlValidationErr error
	loaded            bool
}

// YamlFile - get the yaml file we are working with
func (o *_GlobalOptions) YamlFile() yamlfile.YamlFile { return o.yamlFile }

// IsPipe - whether we are reading the YAML from stdin
func (o *_GlobalOptions) IsPipe() bool { return o.pipe }

// ValidationError - get the YAML validation error (if any)
func (o *_GlobalOptions) ValidationError() error { return o.yamlValidationErr }

// Load - load the yaml file.  This can be called by subcommands for "delayed" parsing
func (o *_GlobalOptions) Load() (err error) {
	if o.loaded {
		return nil
	}
	if o.yamlFile.Exists() {
		if _, err = o.yamlFile.Load(); err != nil {
			if !o.pipe {
				err = errors.Wrapf(err, "File '%s'", o.yamlFile.Filename())
			}

			o.yamlValidationErr = err
		}
	}
	o.loaded = true

	return
}

// NewGoyamlApp - create the goyaml app
func NewGoyamlApp(version, commit, date string) *cli.App {
	rootCmd := newRootCommand(version, commit, date)

	rootCmd.AddSubCommands(
		newGetCommand(rootCmd.GlobalOpts()),
		newSetCommand(rootCmd.GlobalOpts()),
		newDeleteCommand(rootCmd.GlobalOpts()),
		newContainsCommand(rootCmd.GlobalOpts()),
		newValidateCommand(rootCmd.GlobalOpts()),
		newFromJSONCommand(rootCmd.GlobalOpts()),
		newToJSONCommand(rootCmd.GlobalOpts()),
		newExpandCommand(rootCmd.GlobalOpts()),
	)
	return cli.NewApp(rootCmd)
}
