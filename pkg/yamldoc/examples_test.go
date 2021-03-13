package yamldoc

import (
	"bytes"
	"fmt"
	"io"
)

func ExampleNew() {
	// Dummy YAML content
	var yamlStr = `
first:
  second:
    third: a value
`
	// Here creating a reader to some YAML content for this sample code
	var yamlReader io.Reader = bytes.NewBufferString(yamlStr)

	// Create a YamlDoc from the reader
	yamlDoc, err := New(yamlReader)
	if err != nil {
		fmt.Printf("ERROR while parsing yaml: %v\n", err.Error())
		return
	}

	// Read a value from the YAML file
	var value string
	if value, err = yamlDoc.GetString("first.second.third"); err != nil {
		fmt.Printf("ERROR while reading value from yaml\n")
		return
	}

	fmt.Printf("Value read: %s\n", value)
	// Output: Value read: a value
}

func ExampleFromString() {
	// Dummy YAML content
	var yamlStr = `
first:
  second:
    third: a value
`
	// In this example, the yaml is comming from a variable, but all
	// we need is a string with YAML content.
	yamlDoc, err := FromString(yamlStr)
	if err != nil {
		fmt.Printf("ERROR while parsing yaml: %v\n", err.Error())
		return
	}

	// Read a value from the YAML file
	var value string
	if value, err = yamlDoc.GetString("first.second.third"); err != nil {
		fmt.Printf("ERROR while reading value from yaml\n")
		return
	}

	fmt.Printf("Value read: %s\n", value)
	// Output: Value read: a value
}

func ExampleFromBytes() {
	// Dummy YAML content
	var yamlStr = `
first:
  second:
    third: a value
`

	// Get bytes of YAML text.  In this example, we just get them from the string
	yamlBytes := []byte(yamlStr)

	// Create a yamlDoc from the bytes read
	yamlDoc, err := FromBytes(yamlBytes)
	if err != nil {
		fmt.Printf("ERROR while parsing yaml: %v\n", err.Error())
		return
	}

	// Read a value from the YAML file
	var value string
	if value, err = yamlDoc.GetString("first.second.third"); err != nil {
		fmt.Printf("ERROR while reading value from yaml\n")
		return
	}

	fmt.Printf("Value read: %s\n", value)
	// Output: Value read: a value
}
