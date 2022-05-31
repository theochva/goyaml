package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	_flagFile            = "file"
	_flagFileShort       = "f"
	_flagOutput          = "output"
	_flagOutputShort     = "o"
	_flagDetails         = "details"
	_flagDetailsShort    = "d"
	_flagPretty          = "pretty"
	_flagPrettyShort     = "p"
	_flagInput           = "input"
	_flagInputShort      = "i"
	_flagTemplate        = "template"
	_flagTemplateShort   = "t"
	_flagText            = "text"
	_flagExtensions      = "ext"
	_flagExtensionsShort = "e"
	_flagType            = "type"
	_flagTypeShort       = "t"
	_flagStdin           = "stdin"
)

const (
	_CmdOptValidationAware = "CmdOptValidationAware"
	_CmdOptSkipParsing     = "CmdOptSkipParsing"
	_CmdOptValueTrue       = "true"
	_CmdOptValueFalse      = "false"
)

const (
	_FormatText = "text"
	_FormatJSON = "json"
	_FormatYAML = "yaml"
	_FormatHTML = "html"
)

var (
	// ErrUnsupportedOutputFormat - when an unknown output format was specified
	ErrUnsupportedOutputFormat = fmt.Errorf("unsupport output format.  Supported values are: %s", strings.Join(outputFormatValues, ", "))
	outputFormatValues         = []string{_FormatJSON, _FormatYAML}
)

func validateEnumValues(userValue, errPrefix string, validValues []string) error {
	if userValue != "" {
		valid := false
		for _, value := range validValues {
			if value == userValue {
				valid = true
				break
			}
		}
		if !valid {
			return newValidationError("%s. Valid values are: %s", errPrefix, strings.Join(validValues, ", "))
		}
	}
	return nil
}

func marshalValue(value interface{}, outputFormat string) (bytes []byte, err error) {
	switch outputFormat {
	case _FormatYAML:
		bytes, err = yaml.Marshal(value)
	case _FormatJSON:
		bytes, err = marshalToJSON(value, false)
	default:
		err = ErrUnsupportedOutputFormat
	}

	return
}

func marshalToJSON(v interface{}, indent bool) (bytes []byte, err error) {
	if indent {
		if bytes, err = json.MarshalIndent(v, "", "\t"); err != nil {
			return
		}
	} else {
		if bytes, err = json.Marshal(v); err != nil {
			return
		}
	}

	return
}

func convertBytes(bytes []byte, valueType string) (actualValue interface{}, err error) {
	switch valueType {
	case _FormatText:
		actualValue = string(bytes)
	case _FormatYAML:
		err = yaml.Unmarshal(bytes, &actualValue)
	case _FormatJSON:
		err = json.Unmarshal(bytes, &actualValue)
	}

	return
}

func convertFileValue(filename string, valueType string) (actualValue interface{}, err error) {
	var bytes []byte

	if bytes, err = os.ReadFile(filename); err != nil {
		return
	}

	return convertBytes(bytes, valueType)
}
