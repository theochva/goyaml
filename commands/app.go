package commands

import (
	"fmt"
	"os"
)

// SetVersion - set the app version info
func SetVersion(version, commit, date string) {
	globalOpts.rootCmd.Version = fmt.Sprintf("%s [Build date: %s Commit: %s]", version, date, commit)
}

// Execute - run the application
func Execute() {

	if err := globalOpts.rootCmd.Execute(); err != nil {
		if _, isValidationErr := err.(validationError); !isValidationErr {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}
