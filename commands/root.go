package commands

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/theochva/goyaml/yamlfile"
)

type _GlobalOpts struct {
	rootCmd             *cobra.Command
	file                string
	pipe                bool
	yamlFile            yamlfile.YamlFile
	yamlValidationErr   error
	validationAwareCmds map[string]struct{}
	skipParsingCmds     map[string]struct{}
}

var globalOpts = &_GlobalOpts{
	validationAwareCmds: map[string]struct{}{},
	skipParsingCmds:     map[string]struct{}{},
	rootCmd: &cobra.Command{
		Use:                   replaceProgName("$PROG_NAME <global-flags> <command> [options]"),
		Short:                 "Utility to perform simple get/set/delete from yaml files or stdin",
		DisableFlagsInUseLine: false,
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	return nil
		// },
		Long: `Utility to perform simple operations on YAML files: 
  - get/set/delete/check properties to/from YAML content/file
  - Validate YAML content/file
  - Convert to/from YAML/JSON content/file
  - Expand Go templates using YAML as the values file
	
All actions can be performed using either files or stdin/stdout.
	
Primarily intended to be used in scripts or command line.

RC is always 0 unless there was an error while processing.`,
		Example: replaceProgName(`  $PROG_NAME [-f <yaml_file>] <command> [options]
  $PROG_NAME --file <yaml_file> <command> [options]
  $PROG_NAME -f <yaml_file> <command> [options]
  cat foo.yaml | $PROG_NAME <command> [options]`),
	},
}

func init() {
	// Setup options for the global flags
	globalOpts.rootCmd.PersistentPreRunE = globalOpts.processFlags
	globalOpts.rootCmd.PersistentFlags().StringVarP(
		&globalOpts.file,
		_flagFile, _flagFileShort, "",
		"The yaml file to read/write. If not specified it reads from stdin",
	)

	setExamplesAtEndOfHelp(globalOpts.rootCmd)
	setUsageReturnCode(globalOpts.rootCmd, 1)
}

func (o *_GlobalOpts) processFlags(cmd *cobra.Command, args []string) error {
	// Check if help is requested on a command
	if strings.HasPrefix(cmd.Use, "help") {
		return nil
	}
	// Otherwise, we check the global flags
	o.pipe = (o.file == "")
	if o.yamlFile = newYamlFileWrapper(o.file); o.yamlFile != nil {
		if !o.isSkipParsingCommand(cmd) {
			if o.yamlFile.Exists() {
				if _, err := o.yamlFile.Load(); err != nil {
					if o.pipe {
						o.yamlValidationErr = err
					} else {
						o.yamlValidationErr = errors.Wrapf(err, "File '%s'", o.yamlFile.Filename())
					}

					if !o.isValidationErrAwareCommand(cmd) {
						return o.yamlValidationErr
					}
				}
			}
		}
	}

	return nil
}

func (o *_GlobalOpts) addCommand(cmd *cobra.Command, validationErrAware, skipParsing bool) {
	o.rootCmd.AddCommand(cmd)

	if validationErrAware {
		o.validationAwareCmds[cmd.Use] = struct{}{}
	}
	if skipParsing {
		o.skipParsingCmds[cmd.Use] = struct{}{}
	}
}

func (o *_GlobalOpts) isValidationErrAwareCommand(cmd *cobra.Command) bool {
	_, isValidationAware := o.validationAwareCmds[cmd.Use]
	return isValidationAware
}

func (o *_GlobalOpts) isSkipParsingCommand(cmd *cobra.Command) bool {
	_, isSkipParsing := o.skipParsingCmds[cmd.Use]
	return isSkipParsing
}
