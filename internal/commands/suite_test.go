package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/theochva/go-misc/pkg/osext"
	"github.com/theochva/goyaml/internal/commands/cli"
)

var (
	_SampleYAMLExistingKey         = "xmas-fifth-day.calling-birds"
	_SampleYAMLExistingValue       = "four"
	_SampleYAMLExistingValueAsJSON = "\"four\""
	_SampleYAMLExistingValueAsYAML = "four"

	_SampleYAMLExistingBoolKey         = "xmas"
	_SampleYAMLExistingBoolValue       = "true"
	_SampleYAMLExistingBoolValueAsJSON = "true"
	_SampleYAMLExistingBoolValueAsYAML = "true"

	_SampleYAMLExistingIntKey         = "xmas-fifth-day.golden-rings"
	_SampleYAMLExistingIntValue       = "5"
	_SampleYAMLExistingIntValueAsJSON = "5"
	_SampleYAMLExistingIntValueAsYAML = "5"

	_SampleYAMLExistingArrayKey         = "calling-birds"
	_SampleYAMLExistingArrayValue       = "[huey dewey louie fred]"
	_SampleYAMLExistingArrayValueAsJSON = `["huey","dewey","louie","fred"]`
	_SampleYAMLExistingArrayValueAsYAML = strings.TrimSpace(`
- huey
- dewey
- louie
- fred
`)
	_SampleYAMLNonExistingKey = _SampleYAMLExistingKey + "NotThere"
)

var _SampleNonYAML = strings.TrimSpace(`
This is a string which is not parsable as YAML.

It is not parsable as JSON either.
`)

// Sample yaml from: https://www.cloudbees.com/blog/yaml-tutorial-everything-you-need-get-started/
var _SampleYAML = strings.TrimSpace(`
calling-birds:
- huey
- dewey
- louie
- fred
doe: a deer, a female deer
french-hens: 3
pi: 3.14159
ray: a drop of golden sun
xmas: true
xmas-fifth-day:
  calling-birds: four
  french-hens: 3
  golden-rings: 5
  partridges:
    count: 1
    location: a pear tree
  turtle-doves: two
`)
var _SampleJSON = strings.TrimSpace(`
{
	"calling-birds": [
		"huey",
		"dewey",
		"louie",
		"fred"
	],
	"doe": "a deer, a female deer",
	"french-hens": 3,
	"pi": 3.14159,
	"ray": "a drop of golden sun",
	"xmas": true,
	"xmas-fifth-day": {
		"calling-birds": "four",
		"french-hens": 3,
		"golden-rings": 5,
		"partridges": {
			"count": 1,
			"location": "a pear tree"
		},
		"turtle-doves": "two"
	}
}
`)
var _SampleJSONCompact = `{"calling-birds":["huey","dewey","louie","fred"],"doe":"a deer, a female deer","french-hens":3,"pi":3.14159,"ray":"a drop of golden sun","xmas":true,"xmas-fifth-day":{"calling-birds":"four","french-hens":3,"golden-rings":5,"partridges":{"count":1,"location":"a pear tree"},"turtle-doves":"two"}}`

var (
	testApp                    *cli.App
	testJSONFile, testYAMLFile *os.File
	testNonYAMLFile            *os.File
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Test Suite")
}

func createTestApp() *cli.App {
	os.Setenv("GO_TESTING", "true")
	return NewGoyamlApp("0.0.0", "none", "unknown")
}

func getHelpTextForCommand(cmdName string) string {
	rootCmd := testApp.GetRootCommand().GetCliCommand()

	getHelpText := func(cmd *cobra.Command) string {
		buf := bytes.NewBufferString(cmd.Long)
		buf.WriteString("\n\n")
		buf.WriteString(cmd.UsageString())
		return strings.TrimSpace(buf.String())
	}

	// If cmdName == "" means root
	if cmdName == "" {
		return getHelpText(rootCmd)
	}
	subCommands := rootCmd.Commands()

	for _, subCommand := range subCommands {
		key := strings.Split(subCommand.Use, " ")[0]
		if key == cmdName {
			return getHelpText(subCommand)
		}
	}
	return ""
}

func runCommand(input string, args ...string) (string, error) {
	// func runCommand(app *cli.App, input string, args ...string) (string, error) {
	rootCmd := testApp.GetRootCommand().GetCliCommand()

	rootCmd.SetIn(bytes.NewBufferString(input))

	outBuf := bytes.NewBuffer([]byte{})
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)

	if len(args) > 0 {
		rootCmd.SetArgs(args)
	}
	testApp.Execute()
	outputBytes, err := ioutil.ReadAll(outBuf)
	if err != nil {
		return "", err
	}
	output := strings.TrimSpace(string(outputBytes))

	if os.Getenv("SHOW_CMD_OUTPUT") == "true" {
		fmt.Printf("COMMAND LINE: goyaml %s\n\n", strings.Join(args, " "))
		fmt.Printf("OUTPUT:\n%s\n", output)
	}
	return output, nil
}

var _ = BeforeSuite(func() {
	var err error

	testNonYAMLFile, err = osext.CreateTempWithContents("", "test*.txt", []byte(_SampleNonYAML), 0644)
	Expect(err).ToNot(HaveOccurred())
	Expect(testNonYAMLFile).ToNot(BeNil())

	testYAMLFile, err = osext.CreateTempWithContents("", "test*.yaml", []byte(_SampleYAML), 0644)
	Expect(err).ToNot(HaveOccurred())
	Expect(testYAMLFile).ToNot(BeNil())

	testJSONFile, err = osext.CreateTempWithContents("", "test*.json", []byte(_SampleJSON), 0644)
	Expect(err).ToNot(HaveOccurred())
	Expect(testJSONFile).ToNot(BeNil())
})

var _ = AfterSuite(func() {
	if testNonYAMLFile != nil {
		os.Remove(testNonYAMLFile.Name())
	}

	if testYAMLFile != nil {
		os.Remove(testYAMLFile.Name())
	}

	if testJSONFile != nil {
		os.Remove(testJSONFile.Name())
	}
})

var _ = BeforeEach(func() {
	testApp = createTestApp()
})
