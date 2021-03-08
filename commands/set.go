package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type _SetCmd struct {
	valueType   string
	inputFile   string
	readStdin   bool
	valueSource string
}

const (
	_ValueTypeString = "string"
	_ValueTypeInt    = "int"
	_ValueTypeBool   = "bool"
	_ValueTypeJSON   = _FormatJSON
	_ValueTypeYAML   = _FormatYAML

	_ValueSourceArg   = "arg"
	_ValueSourceFile  = "file"
	_ValueSourceStdin = "stdin"
)

var validValueTypes = []string{_ValueTypeString, _ValueTypeInt, _ValueTypeBool, _ValueTypeJSON, _ValueTypeYAML}

func init() {
	globalOpts.addCommand(
		(&_SetCmd{}).createCLICommand(),
		false, // Dont care for yaml validation errors
		false) // Dont skip parsing yaml
}

func (o *_SetCmd) createCLICommand() *cobra.Command {
	validTypesWithOr := strings.Join(validValueTypes, "|")

	var cmd = &cobra.Command{
		Use: replaceProgName(`set <key> <value> [-t|--type %s]
  $PROG_NAME -f|--file <yaml-file> set <key> --stdin [-t|--type %s]
  $PROG_NAME [-f|--file <yaml-file>] set <key> -i|--input <value-file> [-t|--type %s]`, validTypesWithOr, validTypesWithOr, validTypesWithOr),
		DisableFlagsInUseLine: true,
		Aliases:               []string{"s"},
		Short:                 "Set a value in a YAML document",
		Long:                  "Set a value in a YAML document.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires the 'key' for the value to be set")
			} else if len(args) > 2 {
				return fmt.Errorf("too many arguments")
			}
			return nil
		},
		ArgAliases: []string{"key", "value"},
		PreRunE:    o.validateAndPreProcessParams,
		RunE:       o.run,
		Example: replaceProgName(`  Update a YAML file with a "primitive" value specified as a parameter:
    $PROG_NAME -f /tmp/foo.yaml set first.second.strProp "someValue"
    $PROG_NAME -f /tmp/foo.yaml set first.second.intProp 10 -t int
    $PROG_NAME -f /tmp/foo.yaml set first.second.boolProp true -t bool

  Update YAML read from stdin with a "primitive" value specified as a parameter and print result to stdout:
    cat /tmp/foo.yaml | $PROG_NAME set first.second.strProp "someValue"
    cat /tmp/foo.yaml | $PROG_NAME set first.second.intProp 10 -t int
    cat /tmp/foo.yaml | $PROG_NAME set first.second.boolProp true -t bool

  Update a YAML file with a value from a JSON string:
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "{\"prop1\": \"str-value\"}" -t json
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "{\"prop1\": 100}" -t json
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "{\"prop1\": true}" -t json

  Update YAML read from stdin with a JSON value specified as a parameter and print result to stdout:
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "{\"prop1\": \"str-value\"}" -t json
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "{\"prop1\": 100}" -t json
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "{\"prop1\": true}" -t json
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "{\"prop1\": true}" -t json | $PROG_NAME set first.second.fourth 100 -t int

  Update a YAML file with a value from a YAML string:
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "prop1: str-value" -t yaml
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "prop1: 100" -t yaml
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "prop1: true" -t yaml

  Update YAML read from stdin with a YAML value specified as a parameter and print result to stdout:
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "prop1: str-value" -t yaml
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "prop1: 100" -t yaml
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "prop1: true" -t yaml`),
	}

	cmd.Flags().StringVarP(
		&o.valueType,
		_flagType, _flagTypeShort, _ValueTypeString,
		"the value type to set. Valid values are: "+strings.Join(validValueTypes, ", "),
	)
	cmd.Flags().BoolVarP(
		&o.readStdin,
		_flagStdin, "", false,
		"read stdin for the value to set",
	)
	cmd.Flags().StringVarP(
		&o.inputFile,
		_flagInput, _flagInputShort, "",
		"the file containing the value to set",
	)
	return cmd
}

func (o *_SetCmd) validateAndPreProcessParams(cmd *cobra.Command, args []string) error {
	multiSourceErr := fmt.Errorf(
		"Must select only one source of the value to set. It can be specified either via "+
			"the \"value\" argument, the flag '-%s|--%s' or the flag '--%s'", _flagInputShort, _flagInput, _flagStdin)
	if len(args) == 2 {
		o.valueSource = _ValueSourceArg
	}
	if o.inputFile != "" {
		if o.valueSource != "" {
			return multiSourceErr
		}
		o.valueSource = _ValueSourceFile
	}
	if o.readStdin {
		if o.valueSource != "" {
			return multiSourceErr
		}
		o.valueSource = _ValueSourceStdin
	}
	if o.valueSource == "" {
		return fmt.Errorf("Must clearly specify the source of the value to set. It can be specified either via "+
			"the \"value\" argument, the flag '-%s|--%s' or the flag '--%s'", _flagInputShort, _flagInput, _flagStdin)
	} else if o.readStdin && globalOpts.pipe {
		return fmt.Errorf("Cannot use stdin for both the YAML and the value to set")
	}

	return validateEnumValues(o.valueType, "Invalid value type", validValueTypes)
}

func (o *_SetCmd) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key      = args[0]
		value    interface{}
		valueSet bool
	)

	if value, err = o.getValue(args); err != nil {
		return
	}

	if value != nil {
		if valueSet, err = globalOpts.yamlFile.Set(key, value); err != nil {
			return
		}

		if valueSet {
			err = globalOpts.yamlFile.Save()
		}
	}

	if !globalOpts.pipe {
		fmt.Println(valueSet)
	}
	return
}

func (o *_SetCmd) getValue(args []string) (value interface{}, err error) {
	if o.valueSource == _ValueSourceArg {
		return o.parseValue(args[1], o.valueType)
	}

	var bytes []byte

	if o.inputFile != "" {
		bytes, err = os.ReadFile(o.inputFile)
	} else {
		bytes, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		return nil, err
	}
	return o.parseValue(string(bytes), o.valueType)
}

func (o *_SetCmd) parseValue(value string, valueType string) (actualValue interface{}, err error) {
	actualValue = value
	if valueType != _ValueTypeString {
		switch valueType {
		case _ValueTypeInt:
			actualValue, err = strconv.Atoi(strings.TrimSpace(value))
		case _ValueTypeBool:
			actualValue, err = strconv.ParseBool(strings.TrimSpace(value))
		case _ValueTypeJSON:
			err = json.Unmarshal([]byte(value), &actualValue)
		case _ValueTypeYAML:
			err = yaml.Unmarshal([]byte(value), &actualValue)
		}
	}
	return
}
