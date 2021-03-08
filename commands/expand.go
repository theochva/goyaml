package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type _ExpandCmd struct {
	optTemplateFile string
	optExtensions   string

	templateFiles []string
	extensions    []string
	templateText  string
	outputFormat  string
}

var expandFormats = []string{_FormatText, _FormatHTML}

func init() {
	globalOpts.addCommand(
		(&_ExpandCmd{}).createCLICommand(),
		false, // Dont care for yaml validation errors
		false) // Dont skip parsing yaml
}

func (o *_ExpandCmd) createCLICommand() *cobra.Command {
	var cmd = &cobra.Command{
		// Use:     "expand",
		Use: replaceProgName(`expand -t|--template <template-files> [-e|--ext <extensions>] [-o|--output text|html]
  $PROG_NAME expand --text <template-text> [-o|--output text|html]`),
		DisableFlagsInUseLine: true,
		Aliases:               []string{"e"},
		Short:                 "Expand Go templates using the YAML as the values data. The templates are expanded to stdout",
		Long:                  `Expand Go templates using the YAML as the values data. The templates are expanded to stdout.`,
		Args:                  cobra.NoArgs,
		PreRunE:               o.validateAndPreProcessParams,
		RunE:                  o.run,
		Example: replaceProgName(`Simple case, one template file: 

    $PROG_NAME --file /tmp/foo.yaml expand --template /tmp/foo.tmpl
    $PROG_NAME -f /tmp/foo.yaml e -t /tmp/foo.tmpl --format text
    $PROG_NAME -f /tmp/foo.yaml e -t /tmp/foo.tmpl -f html
	
    cat /tmp/foo.yaml | $PROG_NAME e -t /tmp/foo.tmpl --format text
    cat /tmp/foo.yaml | $PROG_NAME e -t /tmp/foo.tmpl -f html

  Simple text expansion:

    $PROG_NAME -f /tmp/foo.yaml e --text "{{.first.second}}" --format text
    $PROG_NAME -f /tmp/foo.yaml e --text "{{<b>.first.second</b>}}" -f html

    cat /tmp/foo.yaml | $PROG_NAME expand --text "{{.first.second}}" --format text
    cat /tmp/foo.yaml | $PROG_NAME expand --text "{{<b>.first.second</b>}}" -f html

  Multiple template and/or folders (searches '/tmp/tmpl' for *.tmpl and *.template):

    $PROG_NAME -f /tmp/foo.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -f html
    $PROG_NAME -f /tmp/foo.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -f text
	
    cat /tmp/foo.yaml | $PROG_NAME expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -f html
    cat /tmp/foo.yaml | $PROG_NAME expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -f text`),
	}

	cmd.Flags().StringVarP(
		&o.optTemplateFile,
		_flagTemplate, _flagTemplateShort, "",
		"the Go templates (comma separated) to expand. Can also include directories",
	)
	cmd.Flags().StringVarP(
		&o.templateText,
		_flagText, "", "",
		"the inline tempalte text to use when expanding",
	)
	cmd.Flags().StringVarP(
		&o.optExtensions,
		_flagExtensions, _flagExtensionsShort, "tmpl",
		"the extensions (comma separated) for template files to look for in directories",
	)
	cmd.Flags().StringVarP(
		&o.outputFormat,
		_flagOutput, _flagOutputShort, "text",
		"the template format to use. Supported formats are: "+strings.Join(expandFormats, ", "),
	)
	return cmd
}

func (o *_ExpandCmd) validateAndPreProcessParams(cmd *cobra.Command, args []string) error {
	o.templateFiles = splitAndTrim(o.optTemplateFile)
	o.extensions = splitAndTrim(o.optExtensions)

	if o.templateText == "" && len(o.templateFiles) == 0 {
		return fmt.Errorf("One of the '-%s|--%s' or '--%s' is required", _flagTemplateShort, _flagTemplate, _flagText)
	} else if o.templateText != "" && len(o.templateFiles) > 0 {
		return fmt.Errorf("Only one of the '-%s|--%s' or '--%s' flags must be specified", _flagTemplateShort, _flagTemplate, _flagText)
	} else if o.templateText != "" && len(o.extensions) > 0 {
		return fmt.Errorf("You cannot specify extensions with the flag '-%s|--%s' when using an inline template", _flagExtensionsShort, _flagExtensions)
	}

	if err := validateEnumValues(o.outputFormat, "Invalid output format specified", expandFormats); err != nil {
		return err
	}

	return nil
}

func (o *_ExpandCmd) run(cmd *cobra.Command, args []string) (err error) {
	var (
		tmpl      templateWrapper
		fileCount = 0
	)

	tmpl = newTemplateWrapper(o.outputFormat)

	if o.templateText != "" {
		if err = tmpl.newTextTemplate(o.templateText); err != nil {
			return
		}
	} else {
		for _, templateFile := range o.templateFiles {
			if !fileOrDirectoryExists(templateFile) {
				return fmt.Errorf("Template file '%s' does not exist", templateFile)
			}
		}

		// Parse templates
		for _, templateFile := range o.templateFiles {
			if isDirectory(templateFile) {
				for _, ext := range o.extensions {
					var (
						filenames []string
						pattern   = fmt.Sprintf("%s/*.%s", templateFile, ext)
					)
					// First check if there are matches for the pattern.
					if filenames, err = filepath.Glob(pattern); err != nil {
						return err
					}

					// If no matching files, then skip over that extension
					if len(filenames) == 0 {
						continue
					}
					fileCount += len(filenames)
					// If found, then parse those files
					if err = tmpl.parseGlob(pattern); err != nil {
						return
					}
				}
			} else {
				if err = tmpl.parseFiles(templateFile); err != nil {
					return
				}
				fileCount++
			}
		}
	}

	// If no files matched, then simply return an error
	if fileCount == 0 && o.templateText == "" {
		return fmt.Errorf("No matching file(s)")
	}
	return tmpl.execute(os.Stdout, globalOpts.yamlFile.Data())
}
