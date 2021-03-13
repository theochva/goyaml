package commands

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/theochva/goyaml/internal/tests"
)

// TestFromJSONCommand - test suite for the from-json command
func TestFromJSONCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Command 'from-json' scenarios", func() {
	var arrayJSON = `["one","two", "three"]`

	When("No params specified", func() {
		It("prints out the help for the 'from-json' command", func() {
			// goyaml from-json --help
			out, err := runCommand("", "from-json", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("from-json")))
		})
	})
	When("Reading array JSON content from STDIN", func() {
		It("prints an error message", func() {
			// cat array.json | goyaml from-json
			out, err := runCommand(arrayJSON, "from-json")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
	})
	When("Reading non JSON content from STDIN", func() {
		It("prints an error message", func() {
			// echo "This is a test" | goyaml from-json
			out, err := runCommand("This is a test", "from-json")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(HavePrefix("Error:"))
		})
	})
	When("No files specified (i.e. reading from STDIN and writing to STDOUT)", func() {
		It("converts JSON to YAML and prints it to stdout", func() {
			// cat file.json | goyaml from-json
			out, err := runCommand(_SampleJSON, "from-json")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(_SampleYAML))
		})
	})
	When("Reading JSON from file and no target file specified", func() {
		It("converts the JSON file to YAML and prints it to stdout", func() {
			// goyaml from-json -i in.json
			out, err := runCommand("", "from-json", "-i", testJSONFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(_SampleYAML))
		})
	})
	When("Reading JSON from file and target file specified", func() {
		var outFile *os.File

		BeforeEach(func() {
			var err error
			outFile, err = tests.CreateTempFile("testout*.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(outFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if outFile != nil {
				os.Remove(outFile.Name())
			}
		})
		It("converts the JSON file to YAML and writes it to the target file", func() {
			// goyaml -f source.yaml from-json -i in.json
			out, err := runCommand("", "-f", outFile.Name(), "from-json", "--input", testJSONFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			var outFileText string
			outFileText, err = tests.ReadFileToString(outFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			Expect(outFileText).To(Equal(_SampleYAML))
		})
	})
	When("Reading JSON from STDIN and target file specified", func() {
		var outFile *os.File

		BeforeEach(func() {
			var err error
			outFile, err = tests.CreateTempFile("testout*.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(outFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if outFile != nil {
				os.Remove(outFile.Name())
			}
		})
		It("converts the JSON content to YAML and writes it to the target file", func() {
			// cat in.json | goyaml -f out.yaml from-json
			out, err := runCommand(_SampleJSON, "-f", outFile.Name(), "from-json")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(BeEmpty())

			var outFileText string
			outFileText, err = tests.ReadFileToString(outFile.Name(), true)
			Expect(err).ToNot(HaveOccurred())
			Expect(outFileText).To(Equal(_SampleYAML))
		})
	})
})
