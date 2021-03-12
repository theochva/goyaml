package commands

import "fmt"

type validationError error

func newValidationError(format string, a ...interface{}) error {
	return validationError(fmt.Errorf(format, a...))
	// return validationError(fmt.Errorf(format+"\n", a...))
}

// IsValidationErr - check if the error is a validation error
func IsValidationErr(err error) bool {
	_, isValidationErr := err.(validationError)
	return isValidationErr
}
