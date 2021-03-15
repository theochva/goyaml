package commands

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestContainsCommand - test suite for the contains command
func TestContainsCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Command 'contains' scenarios", func() {
	When("No params specified", func() {
		It("prints out help for the 'contains' command with --help option", func() {
			out, err := runCommand("", "contains", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("contains")))
		})
		It("prints out help for the 'contains' command", func() {
			out, err := runCommand("", "help", "contains")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("contains")))
		})
	})
	When("Source YAML is coming from STDIN", func() {
		It("prints 'true' when the specified key exists in the YAML content", func() {
			// cat file.yaml | goyaml contains some.key
			out, err := runCommand(_SampleYAML, "contains", _SampleYAMLExistingKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("true"))
		})
		It("prints 'false' when the specified key does not exist in the YAML content", func() {
			// cat file.yaml | goyaml contains some.key
			out, err := runCommand(_SampleYAML, "contains", _SampleYAMLNonExistingKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("false"))
		})
		It("prints an error message when no key specified", func() {
			// cat file.yaml | goyaml contains
			out, err := runCommand(_SampleYAML, "contains")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
	})
	When("Source YAML is specified with '-f' option)", func() {
		It("prints 'true' when the specified key exists in the YAML file", func() {
			// goyaml -f file.yaml contains some.key
			out, err := runCommand("", "-f", testYAMLFile.Name(), "contains", _SampleYAMLExistingKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("true"))
		})
		It("prints 'false' when the specified key does not exist in the YAML file", func() {
			// goyaml -f file.yaml contains some.key
			out, err := runCommand("", "-f", testYAMLFile.Name(), "contains", _SampleYAMLNonExistingKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("false"))
		})
		It("prints an error message when no key specified", func() {
			// cat file.yaml | goyaml contains some.key
			out, err := runCommand("", "-f", testYAMLFile.Name(), "contains")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
	})
})
