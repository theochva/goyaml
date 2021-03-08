package yamlfile

import (
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/theochva/goyaml/yamldoc"
)

const (
	_SampleYaml = `
a:
  b:
    c: value-c
  d:
    e: false
    f: 10
`
)

var (
	yamlText = strings.TrimSpace(_SampleYaml)
	yamlMap  = map[interface{}]interface{}{
		"a": map[interface{}]interface{}{
			"b": map[interface{}]interface{}{
				"c": "value-c",
			},
			"d": map[interface{}]interface{}{
				"e": false,
				"f": 10,
			},
		},
	}
	yamlStringMap = map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value-c",
			},
			"d": map[string]interface{}{
				"e": false,
				"f": 10,
			},
		},
	}
)

func TestYamlFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YamlFile Test Suite")
}

func checkText(yaml yamldoc.YamlDoc, expectedText string) {
	text, err := yaml.Text()
	Expect(err).To(BeNil())
	Expect(text).To(Equal(expectedText))
}

var _ = Describe("YamlFile functions", func() {
	var yamlFile YamlFile

	Context("There is an initial file", func() {
		var file *os.File
		BeforeEach(func() {
			var err error

			file, err = os.CreateTemp("", "test*.yaml")
			Expect(err).To(BeNil())
			Expect(file).ToNot(BeNil())

			err = os.WriteFile(file.Name(), []byte(yamlText), 0666)
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			if file != nil {
				os.Remove(file.Name())
			}
		})
		It("can load the file", func() {
			// create new yaml file
			yamlFile = New(file.Name())

			// check that it is created
			Expect(yamlFile).ToNot(BeNil())

			// check that file exists
			Expect(yamlFile.Exists()).To(BeTrue())

			// load file
			loaded, err := yamlFile.Load()
			// check that it is loaded without errors
			Expect(loaded).To(BeTrue())
			Expect(err).To(BeNil())

			// check the yaml file's text
			checkText(yamlFile, yamlText)
		})
		It("can create/save a file", func() {
			// Remove any previous file
			os.Remove(file.Name())
			// Create new yaml file
			yamlFile = New(file.Name())
			// It should not exist
			Expect(yamlFile.Exists()).To(BeFalse())
			// Set the data
			yamlFile.SetData(yamlMap)

			checkText(yamlFile, yamlText)

			// Save file
			Expect(yamlFile.Save()).To(BeNil())
			// File should exist and filename should be the name of the file
			Expect(yamlFile.Exists()).To(BeTrue())
			Expect(yamlFile.Filename()).To(Equal(file.Name()))

			// check the yaml file's text
			checkText(yamlFile, yamlText)
		})
	})
})
