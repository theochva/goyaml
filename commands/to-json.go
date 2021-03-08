package commands

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type _ToJSONCmd struct {
	pretty     bool
	outputFile string
}

func init() {
	globalOpts.addCommand(
		(&_ToJSONCmd{}).createCLICommand(),
		false, // Dont care for yaml validation errors
		false) // Dont skip parsing yaml
}

func (o *_ToJSONCmd) createCLICommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "to-json [-o|--output <output-json-file>] [-p|--pretty]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"tj", "tojson", "json"},
		Short:                 "Convert YAML to JSON",
		Args:                  cobra.NoArgs,
		RunE:                  o.run,
		Long: `Convert a YAML document to JSON. 

Note:
  Some ordering might be lost in maps and arrays due to the different way
  maps/arrays are implemented in Go. However, the data should all be intact.`,
		Example: replaceProgName(`  $PROG_NAME --file /tmp/foo.yaml to-json --output foo.json
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

	cmd.Flags().BoolVarP(
		&o.pretty,
		_flagPretty, _flagPrettyShort, false,
		"pretty format the json output",
	)
	cmd.Flags().StringVarP(
		&o.outputFile,
		_flagOutput, _flagOutputShort, "",
		"The file to write the JSON output to. If not specified, the output is printed to stdout",
	)
	return cmd
}

func (o *_ToJSONCmd) run(cmd *cobra.Command, args []string) (err error) {
	if converted, mapData := globalOpts.yamlFile.Map(); converted {
		var bytes []byte

		if bytes, err = marshalToJSON(mapData, o.pretty); err != nil {
			return err
		}

		if o.outputFile != "" {
			if err = os.WriteFile(o.outputFile, bytes, 0644); err != nil {
				return errors.Wrapf(err, "Problem writing to '%s'", o.outputFile)
			}
		} else {
			fmt.Println(string(bytes))
		}
		return
	}

	if globalOpts.pipe {
		return fmt.Errorf("Unable to convert YAML from stdin to JSON")
	}
	return fmt.Errorf("Unable to convert YAML file '%s' to JSON", globalOpts.yamlFile.Filename())
}
