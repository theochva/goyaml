package commands

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestGetCommand - test suite for the get command
func TestGetCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Command 'get' scenarios", func() {
	When("No params specified", func() {
		It("prints out help for the 'get' command", func() {

			out, err := runCommand("", "get", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("get")))
		})
	})
	When("Source YAML is comming from STDIN", func() {
		var (
			keys = []string{
				_SampleYAMLExistingKey, _SampleYAMLExistingKey, _SampleYAMLExistingKey,
				_SampleYAMLExistingArrayKey, _SampleYAMLExistingArrayKey, _SampleYAMLExistingArrayKey,
				_SampleYAMLExistingBoolKey, _SampleYAMLExistingBoolKey, _SampleYAMLExistingBoolKey,
				_SampleYAMLExistingIntKey, _SampleYAMLExistingIntKey, _SampleYAMLExistingIntKey}
			values = []string{
				_SampleYAMLExistingValue, _SampleYAMLExistingValue, _SampleYAMLExistingValue,
				_SampleYAMLExistingArrayValue, _SampleYAMLExistingArrayValueAsJSON, _SampleYAMLExistingArrayValueAsYAML,
				_SampleYAMLExistingBoolValue, _SampleYAMLExistingBoolValueAsJSON, _SampleYAMLExistingBoolValueAsYAML,
				_SampleYAMLExistingIntValue, _SampleYAMLExistingIntValueAsYAML, _SampleYAMLExistingIntValueAsYAML}
			outType = []string{
				"", "json", "yaml",
				"", "json", "yaml",
				"", "json", "yaml",
				"", "json", "yaml",
			}
		)
		for index, key := range keys {
			if outType[index] == "" {
				It("prints out the string value for an existing key in the YAML content", func() {
					// cat file.yaml | goyaml get some.existing.key
					out, err := runCommand(_SampleYAML, "get", key)
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(values[index]))
				})
			} else {
				It(fmt.Sprintf("prints out the %s value for an existing key in the YAML content", strings.ToUpper(outType[index])), func() {
					// cat file.yaml | goyaml get some.existing.key -o json
					// cat file.yaml | goyaml get some.existing.key -o yaml
					out, err := runCommand(_SampleYAML, "get", key, "-o", outType[index])
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(values[index]))
				})
			}
		}
		It("prints nothing for a non existing key in the YAML content", func() {
			// cat file.yaml | goyaml get some.non-existing.key
			out, err := runCommand(_SampleYAML, "get", _SampleYAMLNonExistingKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())
		})
		It("prints an error message when no key specified", func() {
			// cat file.yaml | goyaml get
			out, err := runCommand(_SampleYAML, "get")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
	})
	When("Source YAML is specified with '-f' option", func() {
		var (
			keys = []string{
				_SampleYAMLExistingKey, _SampleYAMLExistingKey, _SampleYAMLExistingKey,
				_SampleYAMLExistingArrayKey, _SampleYAMLExistingArrayKey, _SampleYAMLExistingArrayKey,
				_SampleYAMLExistingBoolKey, _SampleYAMLExistingBoolKey, _SampleYAMLExistingBoolKey,
				_SampleYAMLExistingIntKey, _SampleYAMLExistingIntKey, _SampleYAMLExistingIntKey}
			values = []string{
				_SampleYAMLExistingValue, _SampleYAMLExistingValue, _SampleYAMLExistingValue,
				_SampleYAMLExistingArrayValue, _SampleYAMLExistingArrayValueAsJSON, _SampleYAMLExistingArrayValueAsYAML,
				_SampleYAMLExistingBoolValue, _SampleYAMLExistingBoolValueAsJSON, _SampleYAMLExistingBoolValueAsYAML,
				_SampleYAMLExistingIntValue, _SampleYAMLExistingIntValueAsYAML, _SampleYAMLExistingIntValueAsYAML}
			outType = []string{
				"", "json", "yaml",
				"", "json", "yaml",
				"", "json", "yaml",
				"", "json", "yaml",
			}
		)
		for index, key := range keys {
			if outType[index] == "" {
				It("prints out the string value for an existing key in the YAML content", func() {
					// goyaml -f file.yaml get some.existing.key
					out, err := runCommand("", "-f", testYAMLFile.Name(), "get", key)
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(values[index]))
				})
			} else {
				It(fmt.Sprintf("prints out the %s value for an existing key in the YAML content", strings.ToUpper(outType[index])), func() {
					// goyaml -f file.yaml get some.existing.key -o json
					// goyaml -f file.yaml get some.existing.key -o yaml
					out, err := runCommand("", "-f", testYAMLFile.Name(), "get", key, "-o", outType[index])
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(values[index]))
				})
			}
		}
		It("prints nothing for a non existing key in the YAML file", func() {
			// goyaml -f file.yaml get some.non-existing.key
			out, err := runCommand("", "-f", testYAMLFile.Name(), "get", _SampleYAMLNonExistingKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())
		})
		It("prints an error message when no key specified", func() {
			// cat file.yaml | goyaml get
			out, err := runCommand("", "-f", testYAMLFile.Name(), "get")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
	})
})
