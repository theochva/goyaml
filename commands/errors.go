package commands

import "fmt"

type validationError error

func newValidationError(format string, a ...interface{}) error {
	return validationError(fmt.Errorf(format, a...))
	// return validationError(fmt.Errorf(format+"\n", a...))
}
