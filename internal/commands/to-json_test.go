package commands

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/theochva/go-misc/pkg/osext"
)

// TestToJSONCommand - test suite for the to-json command
func TestToJSONCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Command 'to-json' scenarios", func() {
	var twoDocsSample = `[{"one":"first doc"},{"two":"second doc"}]`
	When("No params specified", func() {
		It("prints out help for the 'to-json' command", func() {
			// goyaml to-json --help
			out, err := runCommand("", "to-json", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("to-json")))
		})
	})
	When("Reading YAML from STDIN and no output target file specified", func() {
		It("converts YAML from stdin to JSON compact and prints to stdout", func() {
			// cat file.yaml | goaml to-json
			out, err := runCommand(_SampleYAML, "to-json")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(_SampleJSONCompact))
		})
		It("converts YAML from stdin to (pretty) formatted JSON and prints to stdout", func() {
			// cat file.yaml | goyaml to-json --pretty
			out, err := runCommand(_SampleYAML, "to-json", "--pretty")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(_SampleJSON))
		})
		It("prints an error message when the YAML content is a seq of documents", func() {
			// cat file.yaml | goyaml to-json --pretty
			out, err := runCommand(twoDocsSample, "to-json", "--pretty")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
		It("prints an error message when the input content is invalid", func() {
			// cat file.yaml | goyaml to-json --pretty
			out, err := runCommand("This is a test", "to-json", "--pretty")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
	})
	When("Read YAML from STDIN and output target file specified", func() {
		var outFile *os.File

		BeforeEach(func() {
			var err error
			outFile, err = os.CreateTemp("", "testout*.json")
			Expect(err).ToNot(HaveOccurred())
			Expect(outFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if outFile != nil {
				os.Remove(outFile.Name())
			}
		})
		It("converts the YAML content to compact JSON and writes it to the output target file", func() {
			// cat file.yaml | goyaml to-json -o out.json
			out, err := runCommand(_SampleYAML, "to-json", "-o", outFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			var outFileText string
			outFileText, err = osext.ReadFileAsString(outFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			Expect(outFileText).To(Equal(_SampleJSONCompact))
		})
		It("converts the YAML content to (pretty) formatted JSON and writes it to the output target file", func() {
			// cat file.yaml | goyaml to-json -o out.json --pretty
			out, err := runCommand(_SampleYAML, "to-json", "-o", outFile.Name(), "--pretty")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			var outFileText string
			outFileText, err = osext.ReadFileAsString(outFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			Expect(outFileText).To(Equal(_SampleJSON))
		})
	})

	When("Source YAML file specified and output target file specified", func() {
		var outFile *os.File

		BeforeEach(func() {
			var err error
			outFile, err = os.CreateTemp("", "testout*.json")
			Expect(err).ToNot(HaveOccurred())
			Expect(outFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if outFile != nil {
				os.Remove(outFile.Name())
			}
		})
		It("converts the contents of the source YAML to compact JSON and writes it to the output target file", func() {
			// goyaml -f source.yaml to-json -o out.json
			out, err := runCommand("", "-f", testYAMLFile.Name(), "to-json", "-o", outFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			var outFileText string
			outFileText, err = osext.ReadFileAsString(outFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			Expect(outFileText).To(Equal(_SampleJSONCompact))
		})
		It("converts the contents of the source YAML to (pretty) formatted JSON and writes it to the output target file", func() {
			// goyaml -f source.yaml to-json -o out.json --pretty
			out, err := runCommand("", "-f", testYAMLFile.Name(), "to-json", "-o", outFile.Name(), "--pretty")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			var outFileText string
			outFileText, err = osext.ReadFileAsString(outFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			Expect(outFileText).To(Equal(_SampleJSON))
		})
	})
	When("Source YAML file specified and output target file not specified", func() {
		It("converts the contents of the source YAML to compact JSON and prints JSON to STDOUT", func() {
			// goyaml -f source.yaml to-json
			out, err := runCommand("", "-f", testYAMLFile.Name(), "to-json")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(_SampleJSONCompact))
		})
		It("converts the contents of the source YAML to (pretty) formatted JSON and prints JSON to STDOUT", func() {
			// goyaml -f source.yaml to-json --pretty
			out, err := runCommand("", "-f", testYAMLFile.Name(), "to-json", "--pretty")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(_SampleJSON))
		})
	})
})
