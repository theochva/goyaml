package yamldoc

import (
	"fmt"
	"reflect"
)

type wrongTypeError struct {
	expectedType reflect.Type
	gotType      reflect.Type
}

func (e *wrongTypeError) Error() string {
	var (
		expectedType = "<NIL>"
		gotType      = "<NIL>"
	)
	if e.expectedType != nil {
		expectedType = e.expectedType.String()
	}
	if e.gotType != nil {
		gotType = e.gotType.String()
	}
	return fmt.Sprintf("Expected type '%s' but got '%s'", expectedType, gotType)
}

// IsWrongTypeError - check if the error is a wrong type error.
//
// This type of error will occur when there is a type mismatch in the type of a YAML value with
// respect to what you "requested" or expected the value to be.  For example:
//
// - calling the GetString() function with a key whose value is not actual a string type in the YAML
// content, will raise this type of error.
//
// - Similarly, calling the GetBool() function with a key whose value is not actually a bool type in
// the YAML content, will raise this type of error
func IsWrongTypeError(err error) bool {
	_, isWrongType := err.(*wrongTypeError)
	return isWrongType
}

// type keyNotFoundError struct {
// 	key string
// }

// func (e *keyNotFoundError) Error() string {
// 	return fmt.Sprintf("Key '%s' not found")
// }

// // IsNotFoundError - check if not found
// func IsNotFoundError(err error) bool {
// 	_, isNotFound := err.(*keyNotFoundError)
// 	return isNotFound
// }
