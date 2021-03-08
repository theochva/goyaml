package commands

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

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
