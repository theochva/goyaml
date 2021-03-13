package commands

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/theochva/goyaml/internal/tests"
)

// TestExpandCommand - test suite for the expand command
func TestExpandCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Expand command scenarios", func() {

	// Template variables
	var (
		_SampleTemplateEnvVarName  = "SAMPLE_ENV_VAR"
		_SampleTemplateEnvVarValue = "SomeEnvVarValue"
		_SampleTemplateValues      = strings.TrimSpace(`
four: valueFour
one: valueOne
three: valueThree
two: valueTwo
`)
		_SampleInlineTemplateText = "One: {{.one}}, Two: {{.two}}, Three: {{.three}} and Four: {{.four}}"
		_SampleExpandedInlineText = "One: valueOne, Two: valueTwo, Three: valueThree and Four: valueFour"
		_SampleTemplateText       = strings.TrimSpace(`
BEGIN file: simple.tmpl
Hello. This is {{.one}}
This is {{.two}}

SAMPLE_ENV_VAR environment variable: {{getenv "SAMPLE_ENV_VAR"}}

Using "simple3.tmpl":
{{template "simple3" .}}
Using "simple4.template":
{{template "simple4" .}}
END file: simple.tmpl
`)
		_SampleTemplate3Text = strings.TrimSpace(`
{{ define "simple3" }}
BEGIN file: tmpl/simple3.tmpl
Hello. This is {{.three}}
This is {{.four}}
END file: tmpl/simple3.tmpl
{{end}}	
`)
		_SampleTemplate4Text = strings.TrimSpace(`
{{ define "simple4" }}
BEGIN file: tmpl/simpl4.template
Hello. This is {{.three}}
This is {{.four}}
END file: tmpl/simple4.template
{{end}}
`)
		_SampleExpandedText = strings.TrimSpace(`
BEGIN file: simple.tmpl
Hello. This is valueOne
This is valueTwo

SAMPLE_ENV_VAR environment variable: SomeEnvVarValue

Using "simple3.tmpl":

BEGIN file: tmpl/simple3.tmpl
Hello. This is valueThree
This is valueFour
END file: tmpl/simple3.tmpl

Using "simple4.template":

BEGIN file: tmpl/simpl4.template
Hello. This is valueThree
This is valueFour
END file: tmpl/simple4.template

END file: simple.tmpl
`)
	)

	When("No params specified", func() {
		It("prints out help", func() {
			out, err := runCommand("", "expand", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("expand")))
		})
	})

	Context("Expand a single template with values file", func() {
		var valuesFile, templateFile *os.File
		BeforeEach(func() {
			var err error
			valuesFile, err = tests.CreateTempFileWithContents("sampleValues*.tmpl", _SampleTemplateValues)
			Expect(err).ToNot(HaveOccurred())
			Expect(valuesFile).ToNot(BeNil())
			templateFile, err = tests.CreateTempFileWithContents("sampleTemplate*.tmpl", _SampleInlineTemplateText)
			Expect(err).ToNot(HaveOccurred())
			Expect(valuesFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if valuesFile != nil {
				os.Remove(valuesFile.Name())
			}
			if templateFile != nil {
				os.Remove(templateFile.Name())
			}
		})
		When("A valid values YAML (from STDIN) but invalid or missing template parameters specified (i.e. -t|--template or --text)", func() {
			It("print an error message", func() {
				// cat values.yaml | goyaml expand
				out, err := runCommand(_SampleTemplateValues, "expand")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("print an error message", func() {
				// cat values.yaml | goyaml expand -t template.tmpl --text "template-text"
				out, err := runCommand(_SampleTemplateValues, "expand", "-t", "foo", "--text", "some-template")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})
		When("A valid values YAML (from STDIN) and expanded with a single template", func() {
			for _, format := range []string{"", "html", "text"} {
				var description string
				if format == "" {
					description = "prints out the correct expanded text to STDOUT"
				} else {
					description = fmt.Sprintf("prints out the correct expanded text (as %s) to STDOUT", strings.ToUpper(format))
				}
				It(description, func() {
					// cat values.yaml | goyaml expand -t singleTemplate.tmpl
					// cat values.yaml | goyaml expand -t singleTemplate.tmpl -o html
					// cat values.yaml | goyaml expand -t singleTemplate.tmpl -o text
					var (
						out string
						err error
					)
					if format == "" {
						out, err = runCommand(_SampleTemplateValues, "expand", "-t", templateFile.Name())
					} else {
						out, err = runCommand(_SampleTemplateValues, "expand", "-t", templateFile.Name(), "-o", format)
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(_SampleExpandedInlineText))
				})
			}
		})
		When("A valid values YAML file is specified via the '-f' option but invalid or missing template parameters specified (i.e. -t|--template or --text)", func() {
			It("print an error message since no template parameter specified", func() {
				// goyaml -f values.yaml expand
				out, err := runCommand("", "-f", valuesFile.Name(), "expand")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
			It("print an error message since both --template and --text options were specified", func() {
				// goyaml -f values.yaml expand -t template.tmpl --text "template-text"
				out, err := runCommand("", "-f", valuesFile.Name(), "expand", "-t", "foo", "--text", "some-template")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})

		When("A valid values YAML file is specified via the '-f' option but non-existant template specified", func() {
			It("print an error message since template does not exist", func() {
				// goyaml -f values.yaml expand
				out, err := runCommand("", "-f", valuesFile.Name(), "expand", "-t", valuesFile.Name()+"foo")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})

		When("A valid values YAML file is specified via the '-f' option but no matching templates found", func() {
			It("print an error message since matches for templates files were found", func() {
				// goyaml -f values.yaml expand -t /tmp -e tmplfoo
				out, err := runCommand("", "-f", valuesFile.Name(), "expand", "-t", path.Dir(valuesFile.Name()), "-e", "tmplfoo")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error:"))
			})
		})
		When("A valid values YAML file is specified via the '-f' option and expanded with a single template", func() {
			for _, format := range []string{"", "html", "text"} {
				var description string
				if format == "" {
					description = "prints out the correct expanded text to STDOUT"
				} else {
					description = fmt.Sprintf("prints out the correct expanded text (as %s) to STDOUT", strings.ToUpper(format))
				}
				It(description, func() {
					// goyaml -f values.yaml expand -t singleTemplate.tmpl
					// goyaml -f values.yaml expand -t singleTemplate.tmpl -o html
					// goyaml -f values.yaml expand -t singleTemplate.tmpl -o text
					var (
						out string
						err error
					)
					if format == "" {
						out, err = runCommand("", "-f", valuesFile.Name(), "expand", "-t", templateFile.Name())
					} else {
						out, err = runCommand("", "-f", valuesFile.Name(), "expand", "-t", templateFile.Name(), "-o", format)
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(_SampleExpandedInlineText))
				})
			}
		})
	})
	Context("Expanding inline template with values file", func() {
		var valuesFile *os.File
		BeforeEach(func() {
			var err error
			valuesFile, err = tests.CreateTempFileWithContents("sampleValues*.tmpl", _SampleTemplateValues)
			Expect(err).ToNot(HaveOccurred())
			Expect(valuesFile).ToNot(BeNil())
		})
		AfterEach(func() {
			if valuesFile != nil {
				os.Remove(valuesFile.Name())
			}
		})
		When("A valid values YAML (from STDIN) and expanded with an inline template", func() {
			for _, format := range []string{"", "html", "text"} {
				var description string
				if format == "" {
					description = "prints out the correct expanded text to STDOUT"
				} else {
					description = fmt.Sprintf("prints out the correct expanded text (as %s) to STDOUT", strings.ToUpper(format))
				}
				It(description, func() {
					// cat values.yaml | goyaml expand --text "{{.template}}"
					// cat values.yaml | goyaml expand --text "{{.template}}" -o html
					// cat values.yaml | goyaml expand --text "{{.template}}" -o text
					var (
						out string
						err error
					)
					if format == "" {
						out, err = runCommand(_SampleTemplateValues, "expand", "--text", _SampleInlineTemplateText)
					} else {
						out, err = runCommand(_SampleTemplateValues, "expand", "--text", _SampleInlineTemplateText, "-o", format)
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(_SampleExpandedInlineText))
				})
			}
		})
		When("A valid values YAML file is specified via the '-f' option and expanded with an inline template", func() {
			for _, format := range []string{"", "html", "text"} {
				var description string
				if format == "" {
					description = "prints out the correct expanded text to STDOUT"
				} else {
					description = fmt.Sprintf("prints out the correct expanded text (as %s) to STDOUT", strings.ToUpper(format))
				}
				It(description, func() {
					// goyaml -f values.yaml expand --text "{{.template}}"
					// goyaml -f values.yaml expand --text "{{.template}}" -o html
					// goyaml -f values.yaml expand --text "{{.template}}" -o text
					var (
						out string
						err error
					)
					if format == "" {
						out, err = runCommand("", "-f", valuesFile.Name(), "expand", "--text", _SampleInlineTemplateText)
					} else {
						out, err = runCommand("", "-f", valuesFile.Name(), "expand", "--text", _SampleInlineTemplateText, "-o", format)
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(_SampleExpandedInlineText))
				})
			}
		})
		When("A valid values YAML (from STDIN) and expanded with an invalid inline template", func() {
			It("prints out an error message related to invalid template", func() {
				// cat values.yaml | goyaml expand --text "{{.template}}"
				out, err := runCommand(_SampleTemplateValues, "expand", "--text", "{{if}"+_SampleInlineTemplateText)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: template:"))
			})
			It("prints out an error message because the -e option is specified with a valid inline template", func() {
				// cat values.yaml | goyaml expand --text "{{.template}}" -e tmpl
				out, err := runCommand(_SampleTemplateValues, "expand", "--text", _SampleInlineTemplateText, "-e", "tmpl")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: "))
			})
			It("prints out an error message because the -o option is specified with an invalid output format", func() {
				// cat values.yaml | goyaml expand --text "{{.template}}" -o tmpl
				out, err := runCommand(_SampleTemplateValues, "expand", "--text", _SampleInlineTemplateText, "-o", "tmpl")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: "))
			})
		})

		When("A valid values YAML file is specified via the '-f' option and expanded with an invalid inline template", func() {
			It("prints out an error message related to invalid template", func() {
				// goyaml -f values.yaml expand --text "{{.template}}"
				out, err := runCommand("", "-f", valuesFile.Name(), "expand", "--text", "{{if}"+_SampleInlineTemplateText)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: template:"))
			})
			It("prints out an error message because the -e option is specified with a valid inline template", func() {
				// goyaml -f values.yaml expand --text "{{.template}}" -e tmpl
				out, err := runCommand("", "-f", valuesFile.Name(), "expand", "--text", _SampleInlineTemplateText, "-e", "tmpl")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: "))
			})
			It("prints out an error message because the -o option is specified with an invalid output format", func() {
				// goyaml -f values.yaml expand --text "{{.template}}" -o tmpl
				out, err := runCommand("", "-f", valuesFile.Name(), "expand", "--text", _SampleInlineTemplateText, "-o", "tmpl")
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: "))
			})
		})
		// When("A valid values YAML file is specified via the '-f' option and inline template")
	})
	Context("Expanding file templates with values file", func() {
		var (
			valuesFile, mainTemplateFile, invalidMainTemplateFile, template3File, template4File *os.File
			templatesParamValue, invalidTemplatesParamValue, templateExtParam                   string
		)
		BeforeEach(func() {
			var err error
			valuesFile, err = tests.CreateTempFileWithContents("sampleValues*.tmpl", _SampleTemplateValues)
			Expect(err).ToNot(HaveOccurred())
			Expect(valuesFile).ToNot(BeNil())
			mainTemplateFile, err = tests.CreateTempFileWithContents("sampleTmpl*.tmpl", _SampleTemplateText)
			Expect(err).ToNot(HaveOccurred())
			Expect(mainTemplateFile).ToNot(BeNil())
			invalidMainTemplateFile, err = tests.CreateTempFileWithContents("invalidSampleTmpl*.tmpl", "{{if}"+_SampleTemplateText)
			Expect(err).ToNot(HaveOccurred())
			Expect(invalidMainTemplateFile).ToNot(BeNil())
			template3File, err = tests.CreateTempFileWithContents("sample3Tmpl*.templ", _SampleTemplate3Text)
			Expect(err).ToNot(HaveOccurred())
			Expect(template3File).ToNot(BeNil())
			template4File, err = tests.CreateTempFileWithContents("sample4Tmpl*.template", _SampleTemplate4Text)
			Expect(err).ToNot(HaveOccurred())
			Expect(template4File).ToNot(BeNil())
			os.Setenv(_SampleTemplateEnvVarName, _SampleTemplateEnvVarValue)
			templatesParamValue = fmt.Sprintf("%s,%s", mainTemplateFile.Name(), path.Dir(template3File.Name()))
			invalidTemplatesParamValue = fmt.Sprintf("%s,%s", invalidMainTemplateFile.Name(), path.Dir(template3File.Name()))
			templateExtParam = "templ,template"
		})
		AfterEach(func() {
			os.Unsetenv(_SampleTemplateEnvVarName)
			files := []*os.File{valuesFile, mainTemplateFile, invalidMainTemplateFile, template3File, template4File}
			for _, file := range files {
				if file != nil {
					os.Remove(file.Name())
				}
			}
		})

		When("A valid values YAML (from STDIN) and expanded with valid templates", func() {
			for _, format := range []string{"", "html", "text"} {
				var description string
				if format == "" {
					description = "prints out the correct expanded text to STDOUT"
				} else {
					description = fmt.Sprintf("prints out the correct expanded text (as %s) to STDOUT", strings.ToUpper(format))
				}
				It(description, func() {
					// cat values.yaml | goyaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template"
					// cat values.yaml | goyaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o html
					// cat values.yaml | goyaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o text
					var (
						out string
						err error
					)
					if format == "" {
						out, err = runCommand(_SampleTemplateValues, "expand", "-t", templatesParamValue, "-e", templateExtParam)
					} else {
						out, err = runCommand(_SampleTemplateValues, "expand", "-t", templatesParamValue, "-e", templateExtParam, "-o", format)
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(_SampleExpandedText))
				})
			}
		})
		When("A valid values YAML file is specified via the '-f' option and expanded with valid templates", func() {
			for _, format := range []string{"", "html", "text"} {
				var description string
				if format == "" {
					description = "prints out the correct expanded text to STDOUT"
				} else {
					description = fmt.Sprintf("prints out the correct expanded text (as %s) to STDOUT", strings.ToUpper(format))
				}
				It(description, func() {
					// goyaml -f values.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template"
					// goyaml -f values.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o html
					// goyaml -f values.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o text
					var (
						out string
						err error
					)
					if format == "" {
						out, err = runCommand("", "-f", valuesFile.Name(), "expand", "-t", templatesParamValue, "-e", templateExtParam)
					} else {
						out, err = runCommand("", "-f", valuesFile.Name(), "expand", "-t", templatesParamValue, "-e", templateExtParam, "-o", format)
					}
					Expect(err).ToNot(HaveOccurred())
					Expect(out).To(Equal(_SampleExpandedText))
				})
			}
		})
		When("A valid values YAML (from STDIN) and expanded with invalid templates", func() {
			It("prints out an error message related to invalid template", func() {
				// cat values.yaml | goyaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template"
				out, err := runCommand(_SampleTemplateValues, "expand", "-t", invalidTemplatesParamValue, "-e", templateExtParam)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: template:"))
			})
		})
		When("A valid values YAML file is specified via the '-f' option and expanded with invalid templates", func() {
			It("prints out an error message related to invalid template", func() {
				// goyaml -f values.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template"
				out, err := runCommand("", "-f", valuesFile.Name(), "expand", "-t", invalidTemplatesParamValue, "-e", templateExtParam)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(HavePrefix("Error: template:"))
			})
		})
	})
})
