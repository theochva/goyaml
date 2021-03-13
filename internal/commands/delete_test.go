package commands

import (
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/theochva/goyaml/internal/tests"
)

// TestDeleteCommand - test suite for the delete command
func TestDeleteCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Command 'delete' scenarios", func() {
	var (
		_extraYAMLKey    = "bird-food"
		_SampleExtraYAML = strings.TrimSpace(`
bird-food:
- Black-oil sunflower seed
- White Proso Millet
- Suet cakes
- Nyjer seed
- Cracked corn`)
		//_SampleExtraJSONCompact = `{"bird-food":["Black-oil sunflower seed","White Proso Millet","Suet cakes","Nyjer seed","Cracked corn"]}`
		_SampleExtraJSON = strings.TrimSpace(`
{
	"bird-food": [
		"Black-oil sunflower seed",
		"White Proso Millet",
		"Suet cakes",
		"Nyjer seed",
		"Cracked corn"
	]
}`)
		testExtraYAMLFile, testExtraJSONFile *os.File
	)
	BeforeEach(func() {
		var err error

		testExtraYAMLFile, err = tests.CreateTempFileWithContents("testExtra*.yaml", _SampleExtraYAML)
		Expect(err).ToNot(HaveOccurred())
		Expect(testExtraYAMLFile).ToNot(BeNil())

		testExtraJSONFile, err = tests.CreateTempFileWithContents("testExtra*.json", _SampleExtraJSON)
		Expect(err).ToNot(HaveOccurred())
		Expect(testExtraJSONFile).ToNot(BeNil())
	})
	AfterEach(func() {
		if testExtraYAMLFile != nil {
			os.Remove(testExtraYAMLFile.Name())
		}

		if testExtraJSONFile != nil {
			os.Remove(testExtraJSONFile.Name())
		}
	})
	When("No params specified", func() {
		It("prints out help for the 'delete' command", func() {
			// goyaml delete --help
			out, err := runCommand("", "delete", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("delete")))
		})
	})
	When("Source YAML is comming from STDIN", func() {
		It("prints 'true' when the specified key exists in the YAML content and is deleted", func() {
			// cat file.yaml | goyaml delete some.existing.key
			inputString := strings.Join([]string{_SampleExtraYAML, _SampleYAML}, "\n")
			out, err := runCommand(inputString, "delete", _extraYAMLKey)
			Expect(err).ToNot(HaveOccurred())
			// Output should be minus the extra YML
			Expect(out).To(Equal(_SampleYAML))
		})
		It("prints 'false' when the specified key does not exist in the YAML content", func() {
			// cat file.yaml | goyaml delete some.non-existing.key
			out, err := runCommand(_SampleYAML, "delete", _SampleYAMLNonExistingKey)
			Expect(err).ToNot(HaveOccurred())
			// Output should be the same as the input (untouched)
			Expect(out).To(Equal(_SampleYAML))
		})
		It("prints an error message when no key specified", func() {
			// cat file.yaml | goyaml delete
			out, err := runCommand(_SampleYAML, "delete")
			Expect(err).ToNot(HaveOccurred())
			// Output should have an error
			Expect(out).To(HavePrefix("Error:"))
		})
	})
	When("Source YAML is specified with '-f' option", func() {
		var (
			testDelFile *os.File
			testDelText string
		)
		BeforeEach(func() {
			var err error
			testDelText = strings.Join([]string{_SampleExtraYAML, _SampleYAML}, "\n")
			testDelFile, err = tests.CreateTempFileWithContents("testdel*.yaml", testDelText)
			Expect(err).ToNot(HaveOccurred())
			Expect(testDelFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if testDelFile != nil {
				os.Remove(testDelFile.Name())
			}
		})
		It("prints 'true' when the specified key exists in the YAML file and is deleted", func() {
			// goyaml -f file.yaml delete some.existing.key
			out, err := runCommand("", "-f", testDelFile.Name(), "delete", _extraYAMLKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("true"))
			// Read the updated file's contents
			contents, err := tests.ReadFileToString(testDelFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			// The text should be minus the "extra" yaml
			Expect(contents).To(Equal(_SampleYAML))
		})
		It("prints 'false' when the specified key does not exist in the YAML file", func() {
			// goyaml -f file.yaml delete some.non-existing.key
			out, err := runCommand("", "-f", testDelFile.Name(), "delete", _SampleYAMLNonExistingKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("false"))
			contents, err := tests.ReadFileToString(testDelFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			// Since nothing deleted, the file contents should be same as originally
			Expect(contents).To(Equal(testDelText))
		})
		It("prints an error message when no key specified", func() {
			// goyaml -f file.yaml delete
			out, err := runCommand("", "-f", testDelFile.Name(), "delete")
			Expect(err).ToNot(HaveOccurred())
			// Output should have an error
			Expect(out).To(HavePrefix("Error:"))
		})
	})
})
