package yamldoc

import (
	"fmt"
)

func convert(value map[interface{}]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	for k, v := range value {
		if strKey, ok := k.(string); ok {
			result[strKey] = v
		} else if stringer, ok := k.(fmt.Stringer); ok {
			result[stringer.String()] = v
		} else {
			strKey = fmt.Sprintf("%v", k)
			result[strKey] = v
		}
	}
	return result
}

func convertNested(value interface{}) interface{} {
	switch x := value.(type) {
	case map[interface{}]interface{}:
		mapValue := map[string]interface{}{}
		for k, v := range x {
			mapValue[k.(string)] = convertNested(v)
		}
		return mapValue
	case []interface{}:
		for i, v := range x {
			x[i] = convertNested(v)
		}
	}
	return value
}
