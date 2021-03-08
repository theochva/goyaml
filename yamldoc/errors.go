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
	return fmt.Sprintf("Expected type '%s' but got '%s'", e.expectedType.String(), e.gotType.String())
}

// IsWrongTypeError - check if the error is a wrong type error
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
