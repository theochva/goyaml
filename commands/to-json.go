package commands

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/theochva/goyaml/commands/cli"
)

type _ToJSONCommand struct {
	cli.AppSubCommand

	globalOpts GlobalOptions
	pretty     bool
	outputFile string
}

func newToJSONCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_ToJSONCommand{
		globalOpts: globalOpts,
	}
	cliCmd := &cobra.Command{
		Use:                   "to-json [-o|--output <output-json-file>] [-p|--pretty]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"tj", "tojson", "json"},
		Short:                 "Convert YAML to JSON",
		Args:                  cobra.NoArgs,
		RunE:                  subCmd.run,
		Long: `Convert a YAML document to JSON. 

Note:
  Some ordering might be lost in maps and arrays due to the different way
  maps/arrays are implemented in Go. However, the data should all be intact.`,
		Example: cli.ReplaceProgName(`  $PROG_NAME --file /tmp/foo.yaml to-json --output foo.json
  $PROG_NAME --file /tmp/foo.yaml to-json --output foo.json --pretty
  $PROG_NAME --file /tmp/foo.yaml to-json -o foo.json
  $PROG_NAME --file /tmp/foo.yaml to-json -o foo.json -p
  $PROG_NAME --file /tmp/foo.yaml to-json
  $PROG_NAME --file /tmp/foo.yaml to-json --pretty
  $PROG_NAME --file /tmp/foo.yaml to-json -p

  cat /tmp/foo.yaml | $PROG_NAME to-json -o foo.json
  cat /tmp/foo.yaml | $PROG_NAME to-json -o foo.json --pretty
  cat /tmp/foo.yaml | $PROG_NAME to-json -o foo.json -p
  cat /tmp/foo.yaml | $PROG_NAME to-json -p | jq -r '.first.second'
  cat /tmp/foo.yaml | $PROG_NAME to-json | jq .`),
	}

	cliCmd.Flags().BoolVarP(
		&subCmd.pretty,
		_flagPretty, _flagPrettyShort, false,
		"pretty format the json output",
	)
	cliCmd.Flags().StringVarP(
		&subCmd.outputFile,
		_flagOutput, _flagOutputShort, "",
		"The file to write the JSON output to. If not specified, the output is printed to stdout",
	)

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)
	return subCmd
}

func (c *_ToJSONCommand) run(cmd *cobra.Command, args []string) (err error) {
	if converted, mapData := c.globalOpts.YamlFile().Map(); converted {
		var bytes []byte

		if bytes, err = marshalToJSON(mapData, c.pretty); err != nil {
			return err
		}

		if c.outputFile != "" {
			if err = os.WriteFile(c.outputFile, bytes, 0644); err != nil {
				return errors.Wrapf(err, "Problem writing to '%s'", c.outputFile)
			}
		} else {
			cmd.Println(string(bytes))
		}
		return
	}

	if c.globalOpts.IsPipe() {
		return fmt.Errorf("Unable to convert YAML from stdin to JSON")
	}
	return fmt.Errorf("Unable to convert YAML file '%s' to JSON", c.globalOpts.YamlFile().Filename())
}
