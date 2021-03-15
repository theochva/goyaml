package commands

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/theochva/go-misc/pkg/osext"
)

// TestSetCommand - test suite for the set command
func TestSetCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Set command scenarios", func() {
	var (
		extraSetStringKey        = "a"
		extraSetStringValue      = "value-a"
		extraSetStringResultYaml = fmt.Sprintf("%s: %s", extraSetStringKey, extraSetStringValue)
		extraSetIntKey           = "a"
		extraSetIntValue         = "100"
		extraSetIntResultYaml    = fmt.Sprintf("%s: %s", extraSetIntKey, extraSetIntValue)
		extraSetBoolKey          = "a"
		extraSetBoolValue        = "true"
		extraSetBoolResultYaml   = fmt.Sprintf("%s: %s", extraSetBoolKey, extraSetBoolValue)
	)
	var extraSetArrayResultYaml = strings.TrimSpace(`
a:
  b:
    c:
    - value1
    - value2
    - value3  
`)
	var extraSetArrayKey = "a.b.c"
	var extraSetArrayValueAsJSON = `["value1","value2","value3"]`
	var extraSetArrayValueAsYAML = strings.TrimSpace(`
- value1
- value2
- value3
`)
	When("No params specified", func() {
		It("prints out help", func() {

			out, err := runCommand("", "set", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("set")))
		})
	})
	Context("Source YAML is specified with '-f' option", func() {
		var workFile *os.File
		BeforeEach(func() {
			var err error
			workFile, err = osext.CreateTempWithContents("", "workfile*.yaml", []byte(_SampleYAML), 0644)
			Expect(err).ToNot(HaveOccurred())
			Expect(workFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if workFile != nil {
				os.Remove(workFile.Name())
			}
		})
		When("A valid value is provided as parameter", func() {
			var (
				types                            = []string{"", "string", "int", "bool", "yaml", "json"}
				keys                             = []string{extraSetStringKey, extraSetStringKey, extraSetIntKey, extraSetBoolKey, extraSetArrayKey, extraSetArrayKey}
				values                           = []string{extraSetStringValue, extraSetStringValue, extraSetIntValue, extraSetBoolValue, extraSetArrayValueAsYAML, extraSetArrayValueAsJSON}
				valuesAsYAML                     = []string{extraSetStringResultYaml, extraSetStringResultYaml, extraSetIntResultYaml, extraSetBoolResultYaml, extraSetArrayResultYaml, extraSetArrayResultYaml}
				description, updatedContent, out string
				err                              error
			)
			// goyaml -f file.yaml set key value
			// goyaml -f file.yaml set key str-value -t string
			// goyaml -f file.yaml set key 10 -t int
			// goyaml -f file.yaml set key true -t bool
			// goyaml -f file.yaml set key "one: value" -t yaml
			// goyaml -f file.yaml set key '{"one": "value"}' -t json
			for index, key := range keys {
				if types[index] == "" {
					description = "Set the value as expected in the YAML file"
				} else {
					description = fmt.Sprintf("Set the value as '%s' in the YAML file", types[index])
				}
				It(description, func() {
					if types[index] == "" {
						out, err = runCommand("", "-f", workFile.Name(), "set", key, values[index])
					} else {
						out, err = runCommand("", "-f", workFile.Name(), "set", key, values[index], "-t", types[index])
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal("true"))
					updatedContent, err = osext.ReadFileAsString(workFile.Name(), true)
					Expect(err).ToNot(HaveOccurred())
					expectedContent := strings.Join([]string{valuesAsYAML[index], _SampleYAML}, "\n")
					Expect(updatedContent).To(Equal(expectedContent))
				})
			}
		})
		When("A valid value is specified from a file", func() {
			var (
				types                            = []string{"", "string", "int", "bool", "yaml", "json"}
				keys                             = []string{extraSetStringKey, extraSetStringKey, extraSetIntKey, extraSetBoolKey, extraSetArrayKey, extraSetArrayKey}
				values                           = []string{extraSetStringValue, extraSetStringValue, extraSetIntValue, extraSetBoolValue, extraSetArrayValueAsYAML, extraSetArrayValueAsJSON}
				valuesAsYAML                     = []string{extraSetStringResultYaml, extraSetStringResultYaml, extraSetIntResultYaml, extraSetBoolResultYaml, extraSetArrayResultYaml, extraSetArrayResultYaml}
				description, updatedContent, out string
				filesToDel                       []*os.File
				err                              error
			)
			BeforeEach(func() {
				filesToDel = []*os.File{}
			})
			AfterEach(func() {
				if len(filesToDel) > 0 {
					for _, file := range filesToDel {
						if file != nil {
							os.Remove(file.Name())
						}
					}
				}
			})
			// goyaml -f file.yaml set key -i value.txt
			// goyaml -f file.yaml set key -i value.txt -t string
			// goyaml -f file.yaml set key -i value.txt -t int
			// goyaml -f file.yaml set key -i value.txt -t bool
			// goyaml -f file.yaml set key -i value.yaml -t yaml
			// goyaml -f file.yaml set key -i value.json -t json
			for index, key := range keys {
				if types[index] == "" {
					description = "Set the value as expected in the YAML file"
				} else {
					description = fmt.Sprintf("Set the value as '%s' in the YAML file", types[index])
				}
				It(description, func() {
					valueFile, err2 := osext.CreateTempWithContents("", "valuesTest*.txt", []byte(values[index]), 0644)
					Expect(err2).ToNot(HaveOccurred())
					Expect(valueFile).ToNot(BeNil())
					filesToDel = append(filesToDel, valueFile)

					if types[index] == "" {
						out, err = runCommand("", "-f", workFile.Name(), "set", key, "-i", valueFile.Name())
					} else {
						out, err = runCommand("", "-f", workFile.Name(), "set", key, "-i", valueFile.Name(), "-t", types[index])
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal("true"))
					updatedContent, err = osext.ReadFileAsString(workFile.Name(), true)
					Expect(err).ToNot(HaveOccurred())
					expectedContent := strings.Join([]string{valuesAsYAML[index], _SampleYAML}, "\n")
					Expect(updatedContent).To(Equal(expectedContent))
				})
			}
		})
		When("A valid value is provided from STDIN", func() {
			var (
				types                            = []string{"", "string", "int", "bool", "yaml", "json"}
				keys                             = []string{extraSetStringKey, extraSetStringKey, extraSetIntKey, extraSetBoolKey, extraSetArrayKey, extraSetArrayKey}
				values                           = []string{extraSetStringValue, extraSetStringValue, extraSetIntValue, extraSetBoolValue, extraSetArrayValueAsYAML, extraSetArrayValueAsJSON}
				valuesAsYAML                     = []string{extraSetStringResultYaml, extraSetStringResultYaml, extraSetIntResultYaml, extraSetBoolResultYaml, extraSetArrayResultYaml, extraSetArrayResultYaml}
				description, updatedContent, out string
				err                              error
			)
			// cat value.txt | goyaml -f file.yaml set key --stdin
			// cat value.txt | goyaml -f file.yaml set key --stdin -t string
			// cat value.txt | goyaml -f file.yaml set key --stdin -t int
			// cat value.txt | goyaml -f file.yaml set key --stdin -t bool
			// cat value.yaml | goyaml -f file.yaml set key --stdin -t yaml
			// cat value.json | goyaml -f file.yaml set key --stdin -t json
			for index, key := range keys {
				if types[index] == "" {
					description = "Set the value as expected in the YAML file"
				} else {
					description = fmt.Sprintf("Set the value as '%s' in the YAML file", types[index])
				}
				It(description, func() {
					if types[index] == "" {
						out, err = runCommand(values[index], "-f", workFile.Name(), "set", key, "--stdin")
					} else {
						out, err = runCommand(values[index], "-f", workFile.Name(), "set", key, "--stdin", "-t", types[index])
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal("true"))
					updatedContent, err = osext.ReadFileAsString(workFile.Name(), true)
					Expect(err).ToNot(HaveOccurred())
					expectedContent := strings.Join([]string{valuesAsYAML[index], _SampleYAML}, "\n")
					Expect(updatedContent).To(Equal(expectedContent))
				})
			}
		})
		When("An invalid value is provided as parameter", func() {
			var (
				types  = []string{"int", "bool", "yaml", "json"}
				keys   = []string{extraSetIntKey, extraSetBoolKey, extraSetArrayKey, extraSetArrayKey}
				values = []string{extraSetIntValue + "f", extraSetBoolValue + "f", ":" + extraSetArrayValueAsYAML, "{" + extraSetArrayValueAsJSON}
			)
			// goyaml -f file.yaml set key 10s -t int
			// goyaml -f file.yaml set key trueF -t bool
			// goyaml -f file.yaml set key "one value" -t yaml
			// goyaml -f file.yaml set key '{"one: "value"}' -t json
			for index, key := range keys {
				It("Prints out an error message", func() {
					out, err := runCommand("", "-f", workFile.Name(), "set", key, values[index], "-t", types[index])
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(HavePrefix("Error:"))
				})
			}
		})
		When("Wrong number of parameters are specified", func() {
			It("Prints an error message when neither key nor value specified", func() {
				// goyaml -f file.yaml set
				out, err := runCommand("", "-f", workFile.Name(), "set")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("Prints an error message when there are missing parameters", func() {
				// goyaml -f file.yaml set someKey
				out, err := runCommand("", "-f", workFile.Name(), "set", "someKey")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("Prints out an error message when extra parameters are specified", func() {
				// goyaml -f file.yaml set someKey someValue extraParam
				out, err := runCommand("", "-f", workFile.Name(), "set", "someKey", "someValue", "extraParam")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})
		When("A value is specified from multiple sources", func() {
			It("Prints an error message when value specified as a parameter and also from STDIN", func() {
				// cat value.txt | goyaml -f file.yaml set key --stdin value
				out, err := runCommand(extraSetStringValue, "-f", workFile.Name(), "set", "--stdin", extraSetStringKey, extraSetStringValue)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("Prints an error message when value commin from STDIN and coming from file", func() {
				// cat value.txt | goyaml -f file.yaml set key --stdin value
				out, err := runCommand(extraSetStringValue, "-f", workFile.Name(), "set", "--stdin", extraSetStringKey, "-i", workFile.Name())
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("Prints an error message when value specified as a parameter and coming from file", func() {
				// cat value.txt | goyaml -f file.yaml set key --stdin value
				out, err := runCommand(extraSetStringValue, "-f", workFile.Name(), "set", extraSetStringKey, extraSetStringValue, "-i", workFile.Name())
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})
	})

	Context("Source YAML is provided via STDIN", func() {
		When("A valid value is provided as parameter", func() {
			var (
				types            = []string{"", "string", "int", "bool", "yaml", "json"}
				keys             = []string{extraSetStringKey, extraSetStringKey, extraSetIntKey, extraSetBoolKey, extraSetArrayKey, extraSetArrayKey}
				values           = []string{extraSetStringValue, extraSetStringValue, extraSetIntValue, extraSetBoolValue, extraSetArrayValueAsYAML, extraSetArrayValueAsJSON}
				valuesAsYAML     = []string{extraSetStringResultYaml, extraSetStringResultYaml, extraSetIntResultYaml, extraSetBoolResultYaml, extraSetArrayResultYaml, extraSetArrayResultYaml}
				description, out string
				err              error
			)
			// cat file.yaml | goyaml set key value
			// cat file.yaml | goyaml set key str-value -t string
			// cat file.yaml | goyaml set key 10 -t int
			// cat file.yaml | goyaml set key true -t bool
			// cat file.yaml | goyaml set key "one: value" -t yaml
			// cat file.yaml | goyaml set key '{"one": "value"}' -t json
			for index, key := range keys {
				if types[index] == "" {
					description = "Set the value as expected in the YAML file"
				} else {
					description = fmt.Sprintf("Set the value as '%s' in the YAML file", types[index])
				}
				It(description, func() {
					if types[index] == "" {
						out, err = runCommand(_SampleYAML, "set", key, values[index])
					} else {
						out, err = runCommand(_SampleYAML, "set", key, values[index], "-t", types[index])
					}
					Expect(err).ToNot(HaveOccurred())
					expectedContent := strings.Join([]string{valuesAsYAML[index], _SampleYAML}, "\n")
					Expect(out).To(Equal(expectedContent))
				})
			}
		})
		When("A valid value is specified from a file", func() {
			var (
				types            = []string{"", "string", "int", "bool", "yaml", "json"}
				keys             = []string{extraSetStringKey, extraSetStringKey, extraSetIntKey, extraSetBoolKey, extraSetArrayKey, extraSetArrayKey}
				values           = []string{extraSetStringValue, extraSetStringValue, extraSetIntValue, extraSetBoolValue, extraSetArrayValueAsYAML, extraSetArrayValueAsJSON}
				valuesAsYAML     = []string{extraSetStringResultYaml, extraSetStringResultYaml, extraSetIntResultYaml, extraSetBoolResultYaml, extraSetArrayResultYaml, extraSetArrayResultYaml}
				description, out string
				filesToDel       []*os.File
				err              error
			)
			BeforeEach(func() {
				filesToDel = []*os.File{}
			})
			AfterEach(func() {
				if len(filesToDel) > 0 {
					for _, file := range filesToDel {
						if file != nil {
							os.Remove(file.Name())
						}
					}
				}
			})
			// cat file.yaml | goyaml set key -i value.txt
			// cat file.yaml | goyaml set key -i value.txt -t string
			// cat file.yaml | goyaml set key -i value.txt -t int
			// cat file.yaml | goyaml set key -i value.txt -t bool
			// cat file.yaml | goyaml set key -i value.yaml -t yaml
			// cat file.yaml | goyaml set key -i value.json -t json
			for index, key := range keys {
				if types[index] == "" {
					description = "Set the value as expected in the YAML file"
				} else {
					description = fmt.Sprintf("Set the value as '%s' in the YAML file", types[index])
				}
				It(description, func() {
					valueFile, err2 := osext.CreateTempWithContents("", "valuesTest*.txt", []byte(values[index]), 0644)
					Expect(err2).ToNot(HaveOccurred())
					Expect(valueFile).ToNot(BeNil())
					filesToDel = append(filesToDel, valueFile)

					if types[index] == "" {
						out, err = runCommand(_SampleYAML, "set", key, "-i", valueFile.Name())
					} else {
						out, err = runCommand(_SampleYAML, "set", key, "-i", valueFile.Name(), "-t", types[index])
					}
					Expect(err).ToNot(HaveOccurred())
					expectedContent := strings.Join([]string{valuesAsYAML[index], _SampleYAML}, "\n")
					Expect(out).To(Equal(expectedContent))
				})
			}
		})
		When("An invalid value is provided as parameter", func() {
			var (
				types  = []string{"int", "bool", "yaml", "json"}
				keys   = []string{extraSetIntKey, extraSetBoolKey, extraSetArrayKey, extraSetArrayKey}
				values = []string{extraSetIntValue + "f", extraSetBoolValue + "f", ":" + extraSetArrayValueAsYAML, "{" + extraSetArrayValueAsJSON}
			)
			// cat file.yaml | goyaml set key 10s -t int
			// cat file.yaml | goyaml set key trueF -t bool
			// cat file.yaml | goyaml set key "one value" -t yaml
			// cat file.yaml | goyaml set key '{"one: "value"}' -t json
			for index, key := range keys {
				It("Prints out an error message", func() {
					out, err := runCommand(_SampleYAML, "set", key, values[index], "-t", types[index])
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(HavePrefix("Error:"))
				})
			}
		})
		When("Wrong number of parameters are specified", func() {
			// cat file.yaml | goyaml set
			It("Prints an error message when neither key nor value specified", func() {
				out, err := runCommand(_SampleYAML, "set")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("Prints an error message when there are missing parameters", func() {
				// cat file.yaml | goyaml set someKey
				out, err := runCommand(_SampleYAML, "set", "someKey")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("Prints out an error message when extra parameters are specified", func() {
				// cat file.yaml | goyaml set someKey someValue extraParam
				out, err := runCommand(_SampleYAML, "set", "someKey", "someValue", "extraParam")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})
		When("A value is specified from multiple sources", func() {
			It("Prints an error message when value specified as a parameter and coming from file", func() {
				// cat value.txt | goyaml -f file.yaml set key --stdin value
				out, err := runCommand(_SampleYAML, "set", "--stdin", extraSetStringKey, "-i", "someFile.txt")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})
	})
})
