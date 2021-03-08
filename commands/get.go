package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/theochva/goyaml/commands/cli"
)

type _GetCommand struct {
	cli.AppSubCommand

	globalOpts   GlobalOptions
	outputFormat string
}

// newGetCommand - create the "get" subcommand
func newGetCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_GetCommand{
		globalOpts: globalOpts,
	}

	cliCmd := &cobra.Command{
		Use:                   fmt.Sprintf("get <key> [-o|--output %s]", strings.Join(outputFormatValues, "|")),
		DisableFlagsInUseLine: true,
		Aliases:               []string{"g"},
		Short:                 "Read a value from the yaml",
		Long:                  "Read a value from the yaml.  You can optionally specify the output format for the retrieved value.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires the 'key' to retrieve")
			}
			return nil
		},
		ArgAliases: []string{"key"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateEnumValues(subCmd.outputFormat, "Invalid output format specified", outputFormatValues)
		},
		RunE: subCmd.run,
		Example: cli.ReplaceProgName(`  $PROG_NAME -f /tmp/foo.yaml get first.second.third
  $PROG_NAME -f /tmp/foo.yaml get first.second.third -o json
  $PROG_NAME --file /tmp/foo.yaml get first.second.third
  $PROG_NAME --file /tmp/foo.yaml get first.second.third --output json

  cat /tmp/foo.yaml | $PROG_NAME get first.second.third
  cat /tmp/foo.yaml | $PROG_NAME get first.second.third -o json
  cat /tmp/foo.yaml | $PROG_NAME get first.second.third --output json`),
	}

	cliCmd.Flags().StringVarP(
		&subCmd.outputFormat,
		_flagOutput, _flagOutputShort, "",
		fmt.Sprintf("the output format for value retrieved. Support formats are: %s", strings.Join(outputFormatValues, ", ")))

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)
	return subCmd
}

func (c *_GetCommand) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key   = args[0]
		value interface{}
	)

	if value, err = c.globalOpts.YamlFile().Get(key); err != nil {
		return
	} else if value == nil {
		return
	}

	if c.outputFormat != "" {
		var bytes []byte
		if bytes, err = marshalValue(value, c.outputFormat); err != nil {
			return
		}
		cmd.Println(string(bytes))
		return
	}

	cmd.Println(value)
	return
}
