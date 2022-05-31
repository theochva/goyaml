package utils

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Inspired from:
//	- Helm functions: https://github.com/helm/helm/blob/master/pkg/engine/funcs.go
//	- Sprig functions: http://masterminds.github.io/sprig/

// mustToTOML - convert an object to a TOML string and returns any errors.
func mustToTOML(v interface{}) (string, error) {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// toTOML - marshals an interface to TOML and returns it as string. Errors are not returned
func toTOML(v interface{}) string {
	str, _ := mustToTOML(v)
	return str
}

// mustToYAML - converts an object to a YAML string and returns any errors
func mustToYAML(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(data), "\n"), nil
}

// toYAML - marshals an interface to YAML and returns it as string.  Errors are not returned
func toYAML(v interface{}) string {
	str, _ := mustToYAML(v)
	return str
}

// mustFromYAML - converts a YAML string into a structured value and returns any errors
func mustFromYAML(str string) (val interface{}, err error) {
	err = yaml.Unmarshal([]byte(str), &val)
	if err != nil {
		return nil, err
	}
	return
}

// fromYAML - converts a YAML string into a structured value. Errors are not returned
func fromYAML(str string) interface{} {
	val, _ := mustFromYAML(str)
	return val
}

// mustToJSON - converts an object to a JSON string and returns any errors
func mustToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// toJSON - marshals an interface to JSON and returns it as string.  Errors are not returned
func toJSON(v interface{}) string {
	str, _ := mustToJSON(v)
	return str
}

// mustFromJSON - converts a JSON string into a structured value and returns any errors
func mustFromJSON(str string) (val interface{}, err error) {
	err = json.Unmarshal([]byte(str), &val)
	return
}

// fromJSON - converts a JSON string into a structured value. Errors are not returned
func fromJSON(str string) interface{} {
	val, _ := mustFromJSON(str)
	return val
}
