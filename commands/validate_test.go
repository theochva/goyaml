package commands

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TestValidateCommand - test suite for the validate command
func TestValidateCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Command 'validate' scenarios", func() {
	When("No params specified", func() {
		It("prints out help for the 'contains' command", func() {
			// goyaml validate --help
			out, err := runCommand("", "validate", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("validate")))
		})
	})
	When("Source YAML is comming from STDIN", func() {
		It("outputs 'true' for valid YAML", func() {
			// cat file.yaml | goyaml validate
			out, err := runCommand(_SampleYAML, "validate")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("true"))
		})
		It("outputs nothing for valid YAML with --details flag", func() {
			// cat file.yaml | goyaml validate --details
			out, err := runCommand(_SampleYAML, "validate", "--details")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())
		})
		It("outputs 'false' for invalid YAML", func() {
			// cat file.txt | goyaml validate
			out, err := runCommand(_SampleNonYAML, "validate")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("false"))
		})
		It("outputs a validation msg for invalid YAML", func() {
			// cat file.txt | goyaml validate --details
			out, err := runCommand(_SampleNonYAML, "validate", "--details")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).ToNot(BeEmpty())
		})
	})
	When("Source YAML is specified with '-f' option", func() {
		It("outputs 'true' for valid YAML", func() {
			// goyaml -f file.yaml validate
			out, err := runCommand("", "-f", testYAMLFile.Name(), "validate")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("true"))
		})
		It("outputs nothing for valid YAML with --details flag", func() {
			// goyaml -f file.yaml validate --details
			out, err := runCommand("", "-f", testYAMLFile.Name(), "validate", "--details")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())
		})
		It("outputs 'false' for invalid YAML", func() {
			// goyaml -f file.txt validate
			out, err := runCommand("", "-f", testNonYAMLFile.Name(), "validate")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal("false"))
		})
		It("outputs a validation msg for invalid YAML", func() {
			// goyaml -f file.txt validate --details
			out, err := runCommand("", "-f", testNonYAMLFile.Name(), "validate", "--details")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).ToNot(BeEmpty())
		})
	})
})
