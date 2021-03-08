package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
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
	_FormatText = "text"
	_FormatJSON = "json"
	_FormatYAML = "yaml"
	_FormatHTML = "html"
)

var (
	// ErrUnsupportedOutputFormat - when an unknown output format was specified
	ErrUnsupportedOutputFormat = fmt.Errorf("Unsupport output format.  Supported values are: %s", strings.Join(outputFormatValues, ", "))
	outputFormatValues         = []string{_FormatJSON, _FormatYAML}
)

func getOutputFormat(outputFormat string) string {
	if outputFormat == "" {
		return _FormatYAML
	}
	return outputFormat
}

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
		if bytes, err = json.MarshalIndent(v, "", "    "); err != nil {
			return
		}
	} else {
		if bytes, err = json.Marshal(v); err != nil {
			return
		}
	}

	return
}

func splitAndTrim(valuesStr string) (values []string) {
	if valuesStr == "" {
		return
	}
	currValues := strings.Split(valuesStr, ",")

	for _, value := range currValues {
		if trimmedValue := strings.TrimSpace(value); trimmedValue != "" {
			values = append(values, trimmedValue)
		}
	}

	return
}
