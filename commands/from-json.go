package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type _FromJSONCmd struct {
	inputFile string
}

func init() {
	globalOpts.addCommand(
		(&_FromJSONCmd{}).createCLICommand(),
		false, // Dont care for yaml validation errors
		true)  // Do skip parsing yaml
}

func (o *_FromJSONCmd) createCLICommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "from-json [-i|--input <input-json-file>]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"fj", "fromjson"},
		Short:                 "Convert JSON to YAML",
		Args:                  cobra.NoArgs,
		RunE:                  o.run,
		Long: `Convert a JSON document (either from stdin or a file) to YAML.

Note:
  Some ordering might be lost in maps and arrays due to the different way
  maps/arrays are implemented in Go. However, the data should all be intact.`,
		Example: replaceProgName(`  $PROG_NAME --file /tmp/foo.yaml from-json --input foo.json
  $PROG_NAME --file /tmp/foo.yaml from-json -i foo.json

  Convert JSON file to YAML and print to stdout:
    $PROG_NAME from-json --input foo.json
    $PROG_NAME from-json -i foo.json
	cat foo.json | $PROG_NAME from-json

  Convert JSON from stdin and write to YAML file:
    cat /tmp/foo.json | $PROG_NAME --file /tmp/foo.yaml from-json
    cat /tmp/foo.json | $PROG_NAME -f /tmp/foo.yaml from-json

  Convert JSON from stdin to YAML and print to stdout:
    cat /tmp/foo.json | $PROG_NAME from-json`),
	}

	cmd.Flags().StringVarP(
		&o.inputFile,
		_flagInput, _flagInputShort, "",
		"The input JSON file to convert to YAML",
	)
	return cmd
}

func (o *_FromJSONCmd) run(cmd *cobra.Command, args []string) (err error) {
	var (
		value   interface{}
		changed = false
	)

	if o.inputFile == "" {
		// Read JSON from stdin
		var bytes []byte

		if bytes, err = ioutil.ReadAll(os.Stdin); err != nil {
			return
		}
		if value, err = convertBytes(bytes, _FormatJSON); err != nil {
			return
		}
	} else {
		// Read JSON from file
		if value, err = convertFileValue(o.inputFile, _FormatJSON); err != nil {
			return
		}
	}

	// Now check the type of object read from the JSON file.
	if mapValue, ok := value.(map[interface{}]interface{}); ok {
		globalOpts.yamlFile.SetData(mapValue)
		changed = true
	} else if strMapValue, ok := value.(map[string]interface{}); ok {
		globalOpts.yamlFile.SetData(nil)
		for k, v := range strMapValue {
			globalOpts.yamlFile.Data()[k] = v
		}
		changed = true
	} else if _, ok := value.([]interface{}); ok {
		return fmt.Errorf("Input JSON is a JSON array and not map-based content")
	} else {
		return fmt.Errorf("Input JSON does not contain map-based content")
	}

	// Changed made, then save the yaml file.
	if changed {
		err = globalOpts.yamlFile.Save()
	}
	return
}
