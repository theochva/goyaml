package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type _GetCmd struct {
	outputFormat string
}

func init() {
	globalOpts.addCommand(
		(&_GetCmd{}).createCLICommand(),
		false, // Dont care for yaml validation errors
		false) // Dont skip parsing yaml
}

func (o *_GetCmd) createCLICommand() *cobra.Command {
	var cmd = &cobra.Command{
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
			return validateEnumValues(o.outputFormat, "Invalid output format specified", outputFormatValues)
		},
		RunE: o.run,
		Example: replaceProgName(`  $PROG_NAME -f /tmp/foo.yaml get first.second.third
  $PROG_NAME -f /tmp/foo.yaml get first.second.third -o json
  $PROG_NAME --file /tmp/foo.yaml get first.second.third
  $PROG_NAME --file /tmp/foo.yaml get first.second.third --output json

  cat /tmp/foo.yaml | $PROG_NAME get first.second.third
  cat /tmp/foo.yaml | $PROG_NAME get first.second.third -o json
  cat /tmp/foo.yaml | $PROG_NAME get first.second.third --output json`),
	}

	cmd.Flags().StringVarP(
		&o.outputFormat,
		_flagOutput, _flagOutputShort, "",
		fmt.Sprintf("the output format for value retrieved. Support formats are: %s", strings.Join(outputFormatValues, ", ")))

	return cmd
}

func (o *_GetCmd) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key   = args[0]
		value interface{}
	)

	if value, err = globalOpts.yamlFile.Get(key); err != nil {
		return
	} else if value == nil {
		return
	}

	if o.outputFormat != "" {
		var bytes []byte
		if bytes, err = marshalValue(value, o.outputFormat); err != nil {
			return
		}
		fmt.Println(string(bytes))
		return
	}

	fmt.Println(value)
	return
}
