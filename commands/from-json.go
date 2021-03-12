package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/theochva/goyaml/commands/cli"
)

type _FromJSONCommand struct {
	cli.AppSubCommand
	SkipParsingCommand

	globalOpts GlobalOptions
	inputFile  string
}

func newFromJSONCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_FromJSONCommand{
		globalOpts: globalOpts,
	}

	cliCmd := &cobra.Command{
		Use:                   "from-json [-i|--input <input-json-file>]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"fj", "fromjson"},
		Short:                 "Convert JSON to YAML",
		Args:                  cobra.NoArgs,
		RunE:                  subCmd.run,
		Long: `Convert a JSON document (either from stdin or a file) to YAML.

Note:
  Some ordering might be lost in maps and arrays due to the different way
  maps/arrays are implemented in Go. However, the data should all be intact.`,
		Example: cli.ReplaceProgName(`  Convert an input file (foo.json) to /tmp/foo.yaml:
    $PROG_NAME --file /tmp/foo.yaml from-json --input foo.json
    $PROG_NAME --file /tmp/foo.yaml from-json -i foo.json

  Convert JSON file to YAML and print to stdout:
    $PROG_NAME from-json --input foo.json
    $PROG_NAME from-json -i foo.json

  Convert JSON from stdin and write to YAML file:
    cat /tmp/foo.json | $PROG_NAME --file /tmp/foo.yaml from-json
    cat /tmp/foo.json | $PROG_NAME -f /tmp/foo.yaml from-json

  Convert JSON from stdin to YAML and print to stdout:
    cat /tmp/foo.json | $PROG_NAME from-json`),
	}

	cliCmd.Flags().StringVarP(
		&subCmd.inputFile,
		_flagInput, _flagInputShort, "",
		"The input JSON file to convert to YAML",
	)

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)
	return subCmd
}

// ShouldSkipParsing - implementing this method to indicate that this command wants to take care of the parsing
func (c *_FromJSONCommand) ShouldSkipParsing() bool { return true }

func (c *_FromJSONCommand) run(cmd *cobra.Command, args []string) (err error) {
	var (
		value   interface{}
		changed = false
	)

	if c.inputFile == "" {
		// Read JSON from stdin
		var bytes []byte

		if bytes, err = ioutil.ReadAll(cmd.InOrStdin()); err != nil {
			return
		}

		if value, err = convertBytes(bytes, _FormatJSON); err != nil {
			return
		}
	} else {
		// Read JSON from file
		if value, err = convertFileValue(c.inputFile, _FormatJSON); err != nil {
			return
		}
	}

	// Now check the type of object read from the JSON file.
	if mapValue, ok := value.(map[interface{}]interface{}); ok {
		c.globalOpts.YamlFile().SetData(mapValue)
		changed = true
	} else if strMapValue, ok := value.(map[string]interface{}); ok {
		c.globalOpts.YamlFile().SetData(nil)
		for k, v := range strMapValue {
			c.globalOpts.YamlFile().Data()[k] = v
		}
		changed = true
	} else if _, ok := value.([]interface{}); ok {
		return fmt.Errorf("Input JSON is a JSON array and not map-based content")
	} else {
		return fmt.Errorf("Input JSON does not contain map-based content")
	}

	// Changed made, then save the yaml file.
	if changed {
		err = c.globalOpts.YamlFile().Save()
	}
	return
}
