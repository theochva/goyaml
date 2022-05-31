package yamldoc

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ErrEmptyKey - error generated when empty key specified
var ErrEmptyKey = errors.New("Empty key specified")

type yamlDoc struct {
	data map[interface{}]interface{}
}

// YamlDoc - interface for manipulating yaml file
type YamlDoc interface {
	// Data - get the underlying map
	Data() map[interface{}]interface{}
	// SetData - set the underlying map
	SetData(newData map[interface{}]interface{}) YamlDoc
	// Map - get the data as map of nested maps with string as the key
	Map() (converted bool, mapData map[string]interface{})
	// Get - get the value at key from the yaml
	Get(key string) (value interface{}, err error)
	// GetString - get the string value at key from the yaml
	GetString(key string) (value string, err error)
	// GetInt - get the int value at key from the yaml
	GetInt(key string) (value int, err error)
	// GetBool - get the bool value at key from the yaml
	GetBool(key string) (value bool, err error)
	// GetObject - get a custom object at key.  The value is unmarshalled into the "obj" parameter
	GetObject(key string, obj interface{}) (err error)
	// Set - get a key from the yaml
	Set(key string, value interface{}) (valueSet bool, err error)
	// Delete - delete a key from the yaml
	Delete(key string) (deleted bool, err error)
	// Contains - check if the specified key path is contained within the yaml
	Contains(key string) (contains bool, err error)
	// Bytes - get the yaml file as bytes
	Bytes() ([]byte, error)
	// Text - get the yaml file as text
	Text() (string, error)
}

// New - create new yaml from reader
func New(reader io.Reader) (YamlDoc, error) {
	result := &yamlDoc{
		data: map[interface{}]interface{}{},
	}

	if reader != nil {
		// Create decoder
		decoder := yaml.NewDecoder(reader)

		if err := decoder.Decode(result.data); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// FromBytes - create new yaml from bytes
func FromBytes(yamlBytes []byte) (YamlDoc, error) {
	return New(bytes.NewBuffer(yamlBytes))
}

// FromString - create new yaml from bytes
func FromString(yamlText string) (YamlDoc, error) {
	return New(bytes.NewBuffer([]byte(yamlText)))
}

// Data - get the underlying map
func (y *yamlDoc) Data() map[interface{}]interface{} {
	return y.data
}

// SetData - set the underlying map
func (y *yamlDoc) SetData(newData map[interface{}]interface{}) YamlDoc {
	if newData == nil {
		y.data = map[interface{}]interface{}{}
	} else {
		y.data = newData
	}
	return y
}

// Map - get the data as map of nested maps with string as the key
func (y *yamlDoc) Map() (converted bool, mapData map[string]interface{}) {
	mapData, converted = convertNested(y.data).(map[string]interface{})

	return

}

// Get - get the value at key from the yaml
func (y *yamlDoc) Get(key string) (value interface{}, err error) {
	if key == "" {
		return
	}

	var (
		keys            = strings.Split(key, ".")
		lastIndex   int = len(keys) - 1
		currData    map[interface{}]interface{}
		containsKey bool
	)

	currData = y.data
	for index, key := range keys {
		if value, containsKey = currData[key]; containsKey {
			if value == nil {
				break
			}
			// If at last key, then we need to return it.  If primitive, we are done.
			// Otherwise perhaps the output flag will be used.
			if index == lastIndex {
				break
			}
			if mapValue, ok := value.(map[interface{}]interface{}); ok {
				currData = mapValue
			}
		} else {
			value = nil
			break
		}
	}
	//
	// Check the value before we return it:
	// - If array, then iterate each values and convert map[interface]interface to map[string]interface
	// - If map[interface]interface, then convert to map[string]interface
	// - Otherwise, leave as is
	//
	if value != nil {
		if array, ok := value.([]interface{}); ok {
			for index, arrValue := range array {
				if mapValue, ok := arrValue.(map[interface{}]interface{}); ok {
					array[index] = convert(mapValue)
				}
			}
			value = array
		} else if mapValue, ok := value.(map[interface{}]interface{}); ok {
			value = convert(mapValue)
		}
	}
	return
}

// GetObject - get a custom object at key.  The value is unmarshalled into the "obj" parameter
func (y *yamlDoc) GetObject(key string, obj interface{}) (err error) {
	var (
		value      interface{}
		valueBytes []byte
	)

	if value, err = y.Get(key); err != nil {
		return err
	}
	if valueBytes, err = yaml.Marshal(&value); err != nil {
		return err
	}
	if err = yaml.Unmarshal(valueBytes, obj); err != nil {
		return err
	}
	return nil
}

// GetString - get the string value at key from the yaml
func (y *yamlDoc) GetString(key string) (value string, err error) {
	var (
		obj    interface{}
		isType bool
	)

	if obj, err = y.Get(key); err != nil {
		return "", err
	}

	if value, isType = obj.(string); !isType {
		return "", &wrongTypeError{
			expectedType: reflect.TypeOf(value),
			gotType:      reflect.TypeOf(obj),
		}
		// return "", errors.Wrapf(err, "Value at key '%s' is not a string", key)
	}
	return
}

// GetBool - get the bool value at key from the yaml
func (y *yamlDoc) GetBool(key string) (value bool, err error) {
	var (
		obj    interface{}
		isType bool
	)

	if obj, err = y.Get(key); err != nil {
		return false, err
	}

	if value, isType = obj.(bool); !isType {
		return false, &wrongTypeError{
			expectedType: reflect.TypeOf(value),
			gotType:      reflect.TypeOf(obj),
		}
		// return false, errors.Wrapf(err, "Value at key '%s' is not a bool", key)
	}
	return
}

// GetInt - get the int value at key from the yaml
func (y *yamlDoc) GetInt(key string) (value int, err error) {
	var (
		obj    interface{}
		isType bool
	)

	if obj, err = y.Get(key); err != nil {
		return 0, err
	}

	if value, isType = obj.(int); !isType {
		return 0, &wrongTypeError{
			expectedType: reflect.TypeOf(value),
			gotType:      reflect.TypeOf(obj),
		}
		// return 0, errors.Wrapf(err, "Value at key '%s' is not a int", key)
	}
	return
}

// Set - get a key from the yaml
func (y *yamlDoc) Set(key string, value interface{}) (valueSet bool, err error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	var (
		keys              = strings.Split(key, ".")
		lastIndex     int = len(keys) - 1
		traversedKeys     = []string{}
		currData      map[interface{}]interface{}
		dataValue     interface{}
	)

	currData = y.data
	for index, key := range keys {
		traversedKeys = append(traversedKeys, key)

		if index == lastIndex {
			currData[key] = value
			valueSet = true
			break
		}

		if dataValue, _ = currData[key]; dataValue == nil {
			dataValue = map[interface{}]interface{}{}
			currData[key] = dataValue
		}
		if mapValue, ok := dataValue.(map[interface{}]interface{}); ok {
			currData = mapValue
		} else {
			return false, fmt.Errorf("key '%s' is not a map container", strings.Join(traversedKeys, "."))
		}
	}
	return
}

// Delete - delete a key from the yaml
func (y *yamlDoc) Delete(key string) (deleted bool, err error) {
	if key == "" {
		return
	}

	var (
		keys          = strings.Split(key, ".")
		lastIndex int = len(keys) - 1
		currData  map[interface{}]interface{}
	)

	currData = y.data
	for index, key := range keys {
		if value, containsKey := currData[key]; containsKey {
			if index == lastIndex {
				delete(currData, key)
				deleted = true
				break
			}
			if mapValue, ok := value.(map[interface{}]interface{}); ok {
				currData = mapValue
			}
		} else {
			deleted = false
			break
		}
	}
	return
}

// Contains - check if the specified key path is contained within the yaml
func (y *yamlDoc) Contains(key string) (contains bool, err error) {
	if key == "" {
		return
	}
	var (
		keys          = strings.Split(key, ".")
		lastIndex int = len(keys) - 1
		currData  map[interface{}]interface{}
	)

	currData = y.data
	for index, key := range keys {
		if value, containsKey := currData[key]; containsKey {
			if index == lastIndex {
				contains = true
				break
			}
			if mapValue, ok := value.(map[interface{}]interface{}); ok {
				currData = mapValue
			}
		} else {
			contains = false
			break
		}
	}
	return
}

// Bytes - get the yaml file as bytes
func (y *yamlDoc) Bytes() ([]byte, error) {
	return yaml.Marshal(y.data)
}

// Text - get the yaml file as text
func (y *yamlDoc) Text() (string, error) {
	bytes, err := y.Bytes()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}
