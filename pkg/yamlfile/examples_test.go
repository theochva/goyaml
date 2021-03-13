package yamlfile

import (
	"fmt"
	"os"
)

func ExampleNew() {
	yamlText := `
first:
  second:
    third: a value
`
	// Create a tmp file to store YAML content
	tmpFile, err := os.CreateTemp("", "testfile*.yaml")
	if err != nil {
		fmt.Printf("ERROR creating temp file: %v\n", err.Error())
		return
	}
	defer os.Remove(tmpFile.Name())

	// Populate the file with some YAML content
	if err = os.WriteFile(tmpFile.Name(), []byte(yamlText), 0644); err != nil {
		fmt.Printf("ERROR while writing to temp file '%s': %v\n", tmpFile.Name(), err.Error())
		return
	}

	// Create an empty YamlFile object
	yamlFile := New(tmpFile.Name())

	if loaded, err := yamlFile.Load(); err != nil {
		fmt.Printf("ERROR while loading YAML file '%s': %v\n", yamlFile.Filename(), err.Error())
		return
	} else if loaded {
		var value string
		if value, err = yamlFile.GetString("first.second.third"); err != nil {
			fmt.Printf("ERROR while reading value from yaml\n")
			return
		}

		fmt.Printf("Value read: %s\n", value)
	}
	// Output: Value read: a value
}

func ExampleLoad() {
	yamlText := `
first:
  second:
    third: a value
`
	// Create a tmp file to store YAML content
	tmpFile, err := os.CreateTemp("", "testfile*.yaml")
	if err != nil {
		fmt.Printf("ERROR creating temp file: %v\n", err.Error())
		return
	}
	defer os.Remove(tmpFile.Name())

	// Populate the file with some YAML content
	if err = os.WriteFile(tmpFile.Name(), []byte(yamlText), 0644); err != nil {
		fmt.Printf("ERROR while writing to temp file '%s': %v\n", tmpFile.Name(), err.Error())
		return
	}

	if loaded, yamlFile, err := Load(tmpFile.Name()); err != nil {
		fmt.Printf("ERROR while loading YAML file '%s': %v\n", yamlFile.Filename(), err.Error())
		return
	} else if loaded {
		var value string
		if value, err = yamlFile.GetString("first.second.third"); err != nil {
			fmt.Printf("ERROR while reading value from yaml\n")
			return
		}

		fmt.Printf("Value read: %s\n", value)
	}
	// Output: Value read: a value
}
