package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/theochva/goyaml/commands/cli"
	"github.com/theochva/goyaml/commands/internal"
)

// _GoyamlRootCommand - the root command for the app
type _GoyamlRootCommand struct {
	cli.AppRootCommand

	globalOpts          *_GlobalOptions
	file                string
	validationAwareCmds map[string]struct{}
	skipParsingCmds     map[string]struct{}
}

// NewRootCommand - create root command
func newRootCommand(version, commit, date string) *_GoyamlRootCommand {
	cliCmd := &cobra.Command{
		Use:                   cli.ReplaceProgName("$PROG_NAME <global-flags> <command> [options]"),
		Short:                 "Utility to perform simple get/set/delete from yaml files or stdin",
		DisableFlagsInUseLine: false,
		Version:               fmt.Sprintf("%s [Build date: %s Commit: %s]", version, date, commit),
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
		Example: cli.ReplaceProgName(`  $PROG_NAME [-f <yaml_file>] <command> [options]
  $PROG_NAME --file <yaml_file> <command> [options]
  $PROG_NAME -f <yaml_file> <command> [options]
  cat foo.yaml | $PROG_NAME <command> [options]`),
	}

	rootCmd := &_GoyamlRootCommand{
		AppRootCommand:      cli.NewAppRootCommandBase(cliCmd),
		file:                "",
		globalOpts:          &_GlobalOptions{},
		validationAwareCmds: map[string]struct{}{},
		skipParsingCmds:     map[string]struct{}{},
	}
	// Setup options for the global flags
	cliCmd.PersistentPreRunE = rootCmd.processFlags
	cliCmd.PersistentFlags().StringVarP(
		&rootCmd.file,
		_flagFile, _flagFileShort, "",
		"The yaml file to read/write. If not specified it reads from stdin",
	)

	// cli.SetVersionWithAuthor(cliCmd, "") //"Bill Theocharoulas - theochva@gmail.com")
	cli.SetExamplesAtEndOfUsage(cliCmd)

	if os.Getenv("GO_TESTING") != "true" {
		cli.SetUsageReturnCode(cliCmd, 1)
	}
	// cli.SetUsageReturnCode(cliCmd, 1)

	return rootCmd
}

// GlobalOpts - get the global opts for the app
func (c *_GoyamlRootCommand) GlobalOpts() GlobalOptions { return c.globalOpts }

func getCommandKey(cmd *cobra.Command) string {
	return strings.Split(cmd.Use, " ")[0]
}

// AddSubCommands - overriding the base root command func to perform some extra things
func (c *_GoyamlRootCommand) AddSubCommands(appSubCmds ...cli.AppSubCommand) {
	// Call parent to add the sub commands
	c.AppRootCommand.AddSubCommands(appSubCmds...)

	// Process the commands and look whether commands are validation aware or should skip parsing
	for _, appSubCmd := range appSubCmds {
		// cliCmd := appSubCmd.GetCliCommand()
		key := getCommandKey(appSubCmd.GetCliCommand())

		if errAware, isAware := appSubCmd.(ValidationErrorAwareCommand); isAware {
			if errAware.IsValidationAware() {
				c.validationAwareCmds[key] = struct{}{}
			}
		}
		if skip, isSkipAware := appSubCmd.(SkipParsingCommand); isSkipAware {
			if skip.ShouldSkipParsing() {
				c.skipParsingCmds[key] = struct{}{}
			}
		}
	}
}

func (c *_GoyamlRootCommand) processFlags(cmd *cobra.Command, args []string) error {
	// Check if help is requested on a command
	if strings.HasPrefix(cmd.Use, "help") {
		return nil
	}
	// Otherwise, we check the global flags
	c.globalOpts.pipe = (c.file == "")
	if c.globalOpts.yamlFile = internal.NewYamlFileWrapper(c.file, cmd.InOrStdin(), cmd.OutOrStdout()); c.globalOpts.yamlFile != nil {
		if !c.isSkipParsingCommand(cmd) {
			if err := c.globalOpts.Load(); err != nil {
				if !c.isValidationErrAwareCommand(cmd) {
					return c.globalOpts.yamlValidationErr
				}
				// if c.globalOpts.yamlFile.Exists() {
				// 	if _, err := c.globalOpts.yamlFile.Load(); err != nil {
				// 		if c.globalOpts.pipe {
				// 			c.globalOpts.yamlValidationErr = err
				// 		} else {
				// 			c.globalOpts.yamlValidationErr = errors.Wrapf(err, "File '%s'", c.globalOpts.yamlFile.Filename())
				// 		}

				// 		if !c.isValidationErrAwareCommand(cmd) {
				// 			return c.globalOpts.yamlValidationErr
				// 		}
				// 	}
			}
		}
	}

	return nil
}

func (c *_GoyamlRootCommand) isValidationErrAwareCommand(cmd *cobra.Command) bool {
	_, isValidationAware := c.validationAwareCmds[getCommandKey(cmd)]
	return isValidationAware
}

func (c *_GoyamlRootCommand) isSkipParsingCommand(cmd *cobra.Command) bool {
	_, isSkipParsing := c.skipParsingCmds[getCommandKey(cmd)]
	return isSkipParsing
}
