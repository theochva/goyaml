package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/theochva/goyaml/internal/commands/cli"
	"gopkg.in/yaml.v2"
)

type _SetCommand struct {
	cli.AppSubCommand

	globalOpts  GlobalOptions
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

// newSetCommand - create the "set" subcommand
func newSetCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_SetCommand{
		globalOpts: globalOpts,
	}

	validTypesWithOr := strings.Join(validValueTypes, "|")
	cliCmd := &cobra.Command{
		Use: cli.ReplaceProgName(`set <key> <value> [-t|--type %s]
  $PROG_NAME -f|--file <yaml-file> set <key> --stdin [-t|--type %s]
  $PROG_NAME [-f|--file <yaml-file>] set <key> -i|--input <value-file> [-t|--type %s]`, validTypesWithOr, validTypesWithOr, validTypesWithOr),
		DisableFlagsInUseLine: true,
		Annotations: map[string]string{
			_CmdOptValidationAware: _CmdOptValueTrue,
			_CmdOptSkipParsing:     _CmdOptValueTrue,
		},
		Aliases: []string{"s"},
		Short:   "Set a value in a YAML document",
		Long: `Set a value in a YAML document. There are multiple ways you can set values in a YAML document:
  - Set a value in a YAML file with a value specified, read from a file or read from stdin  
  - Update a value in a YAML document read from stdin with a value specified or read from a file and print result to stdout`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires the 'key' for the value to be set")
			} else if len(args) > 2 {
				return fmt.Errorf("too many arguments")
			}
			return nil
		},
		ArgAliases: []string{"key", "value"},
		PreRunE:    subCmd.validateAndPreProcessParams,
		RunE:       subCmd.run,
		Example: cli.ReplaceProgName(`  Update a YAML file with a value specified as a parameter:
    $PROG_NAME -f /tmp/foo.yaml set first.second.strProp "someValue"
    $PROG_NAME -f /tmp/foo.yaml set first.second.intProp 10 -t int
    $PROG_NAME -f /tmp/foo.yaml set first.second.boolProp true -t bool
    $PROG_NAME -f /tmp/foo.yaml set first.second.third '{"prop1": "str-value"}' -t json
    $PROG_NAME -f /tmp/foo.yaml set first.second.third '{"prop1": 100}' -t json
    $PROG_NAME -f /tmp/foo.yaml set first.second.third '{"prop1": true}' -t json
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "prop1: str-value" -t yaml
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "prop1: 100" -t yaml
    $PROG_NAME -f /tmp/foo.yaml set first.second.third "prop1: true" -t yaml

  Update a YAML file with a value read from another file:
    $PROG_NAME -f /tmp/foo.yaml set first.second.third -i /tmp/foo.json -t json
    $PROG_NAME -f /tmp/foo.yaml set first.second.third -i /tmp/bar.yaml -t yaml
    $PROG_NAME -f /tmp/foo.yaml set first.second.privateKey -i .ssh/id_rsa_priv

  Update YAML read from stdin with a value specified as a parameter and print result to stdout:
    cat /tmp/foo.yaml | $PROG_NAME set first.second.strProp "someValue"
    cat /tmp/foo.yaml | $PROG_NAME set first.second.intProp 10 -t int
    cat /tmp/foo.yaml | $PROG_NAME set first.second.boolProp true -t bool
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third '{"prop1": "str-value"}' -t json
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third '{"prop1": 100}' -t json
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third '{"prop1": true}' -t json | $PROG_NAME set first.second.fourth 100 -t int
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "prop1: str-value" -t yaml
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "prop1: 100" -t yaml
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third "prop1: true" -t yaml

  Update YAML read from stdin with a value read from another file and print result to stdout:
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third -i /tmp/foo.json -t json
    cat /tmp/foo.yaml | $PROG_NAME set first.second.third -i /tmp/bar.yaml -t yaml
    cat /tmp/foo.yaml | $PROG_NAME set first.second.privateKey -i .ssh/id_rsa_priv

  Generate YAML to stdout with a value read from stdin:
    cat /tmp/foo.json | $PROG_NAME set first.second.third --stdin -t json
    cat /tmp/bar.yaml | $PROG_NAME set first.second.third --stdin -t yaml
    cat ~/.ssh/id_rsa_priv | $PROG_NAME set first.second.privateKey --stdin`),
	}

	cliCmd.Flags().StringVarP(
		&subCmd.valueType,
		_flagType, _flagTypeShort, _ValueTypeString,
		"the value type to set. Valid values are: "+strings.Join(validValueTypes, ", "),
	)
	cliCmd.Flags().BoolVarP(
		&subCmd.readStdin,
		_flagStdin, "", false,
		"read stdin for the value to set",
	)
	cliCmd.Flags().StringVarP(
		&subCmd.inputFile,
		_flagInput, _flagInputShort, "",
		"the file containing the value to set",
	)

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)
	return subCmd
}

func (c *_SetCommand) validateAndPreProcessParams(cmd *cobra.Command, args []string) error {
	// First check if file specified with -f or if value to set is not comming from stdin
	if !c.globalOpts.IsPipe() || !c.readStdin {
		if err := c.globalOpts.Load(); err != nil {
			return err
		}
	}
	multiSourceErr := fmt.Errorf(
		"Must select only one source of the value to set. It can be specified either via "+
			"the \"value\" argument, the flag '-%s|--%s' or the flag '--%s'", _flagInputShort, _flagInput, _flagStdin)
	if len(args) == 2 {
		c.valueSource = _ValueSourceArg
	}
	if c.inputFile != "" {
		if c.valueSource != "" {
			return multiSourceErr
		}
		c.valueSource = _ValueSourceFile
	}
	if c.readStdin {
		if c.valueSource != "" {
			return multiSourceErr
			// } else if c.globalOpts.IsPipe() {
			// 	return fmt.Errorf("Cannot use stdin for both the YAML and the value to set")
		}
		c.valueSource = _ValueSourceStdin
	}
	if c.valueSource == "" {
		return fmt.Errorf("Must clearly specify the source of the value to set. It can be specified either via "+
			"the \"value\" argument, the flag '-%s|--%s' or the flag '--%s'", _flagInputShort, _flagInput, _flagStdin)
	}

	return validateEnumValues(c.valueType, "Invalid value type", validValueTypes)
}

func (c *_SetCommand) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key      = args[0]
		value    interface{}
		valueSet bool
	)

	if value, err = c.getValue(args); err != nil {
		return
	}

	if value != nil {
		if valueSet, err = c.globalOpts.YamlFile().Set(key, value); err != nil {
			return
		}

		if valueSet {
			err = c.globalOpts.YamlFile().Save()
		}
	}

	if !c.globalOpts.IsPipe() {
		cmd.Println(valueSet)
	}
	return
}

func (c *_SetCommand) getValue(args []string) (value interface{}, err error) {
	if c.valueSource == _ValueSourceArg {
		return c.parseValue(args[1], c.valueType)
	}

	var bytes []byte

	if c.inputFile != "" {
		bytes, err = os.ReadFile(c.inputFile)
	} else {
		bytes, err = ioutil.ReadAll(c.GetCliCommand().InOrStdin())
	}
	if err != nil {
		return nil, err
	}
	return c.parseValue(string(bytes), c.valueType)
}

func (c *_SetCommand) parseValue(value string, valueType string) (actualValue interface{}, err error) {
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
