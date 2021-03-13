package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/theochva/goyaml/internal/commands/cli"
	"github.com/theochva/goyaml/internal/commands/utils"
)

type _ExpandCommand struct {
	cli.AppSubCommand

	globalOpts GlobalOptions

	optTemplateFile string
	optExtensions   string

	templateFiles []string
	extensions    []string
	templateText  string
	outputFormat  string
}

var expandFormats = []string{_FormatText, _FormatHTML}

// newGetSubCommand - create the "get" subcommand
func newExpandCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_ExpandCommand{
		globalOpts: globalOpts,
	}

	cliCmd := &cobra.Command{
		Use: cli.ReplaceProgName(`expand -t|--template <template-files> [-e|--ext <extensions>] [-o|--output text|html]
  $PROG_NAME expand --text <template-text> [-o|--output text|html]`),
		DisableFlagsInUseLine: true,
		Aliases:               []string{"e"},
		Short:                 "Expand Go templates using the YAML as the values data. The templates are expanded to stdout",
		Long:                  `Expand Go templates using the YAML as the values data. The templates are expanded to stdout.`,
		Args:                  cobra.NoArgs,
		PreRunE:               subCmd.validateAndPreProcessParams,
		RunE:                  subCmd.run,
		Example: cli.ReplaceProgName(`  Simple case, one template file: 
    $PROG_NAME --file /tmp/foo.yaml expand --template /tmp/foo.tmpl
    $PROG_NAME -f /tmp/foo.yaml e -t /tmp/foo.tmpl --output text
    $PROG_NAME -f /tmp/foo.yaml e -t /tmp/foo.tmpl -o html
	
    cat /tmp/foo.yaml | $PROG_NAME e -t /tmp/foo.tmpl --output text
    cat /tmp/foo.yaml | $PROG_NAME e -t /tmp/foo.tmpl -o html

  Simple text expansion:
    $PROG_NAME -f /tmp/foo.yaml e --text "{{.first.second}}" --output text
    $PROG_NAME -f /tmp/foo.yaml e --text "{{<b>.first.second</b>}}" -o html

    cat /tmp/foo.yaml | $PROG_NAME expand --text "{{.first.second}}" --output text
    cat /tmp/foo.yaml | $PROG_NAME expand --text "{{<b>.first.second</b>}}" -o html

  Multiple template and/or folders (searches '/tmp/tmpl' for *.tmpl and *.template):
    $PROG_NAME -f /tmp/foo.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o html
    $PROG_NAME -f /tmp/foo.yaml expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o text
	
    cat /tmp/foo.yaml | $PROG_NAME expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o html
    cat /tmp/foo.yaml | $PROG_NAME expand -t /tmp/simple.tmpl,/tmp/someTmpl,/tmp/moreTmpl -e "tmpl,template" -o text`),
	}

	cliCmd.Flags().StringVarP(
		&subCmd.optTemplateFile,
		_flagTemplate, _flagTemplateShort, "",
		"the Go templates (comma separated) to expand. Can also include directories",
	)
	cliCmd.Flags().StringVarP(
		&subCmd.templateText,
		_flagText, "", "",
		"the inline tempalte text to use when expanding",
	)
	cliCmd.Flags().StringVarP(
		&subCmd.optExtensions,
		_flagExtensions, _flagExtensionsShort, "tmpl",
		"the extensions (comma separated) for template files to look for in directories",
	)
	cliCmd.Flags().StringVarP(
		&subCmd.outputFormat,
		_flagOutput, _flagOutputShort, "text",
		"the template format to use. Supported formats are: "+strings.Join(expandFormats, ", "),
	)

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)

	return subCmd
}

func (c *_ExpandCommand) validateAndPreProcessParams(cmd *cobra.Command, args []string) error {
	c.templateFiles = splitAndTrim(c.optTemplateFile)
	c.extensions = splitAndTrim(c.optExtensions)

	if c.templateText == "" && len(c.templateFiles) == 0 {
		return fmt.Errorf("One of the '-%s|--%s' or '--%s' is required", _flagTemplateShort, _flagTemplate, _flagText)
	} else if c.templateText != "" && len(c.templateFiles) > 0 {
		return fmt.Errorf("Only one of the '-%s|--%s' or '--%s' flags must be specified", _flagTemplateShort, _flagTemplate, _flagText)
	} else if c.templateText != "" && len(c.extensions) > 0 {
		// If user specified the "-e" flag with the "--text" then this is error
		if flag := cmd.Flag(_flagExtensions); flag != nil && flag.Changed {
			return fmt.Errorf("You cannot specify extensions with the flag '-%s|--%s' when using an inline template", _flagExtensionsShort, _flagExtensions)
		}
	}

	if err := validateEnumValues(c.outputFormat, "Invalid output format specified", expandFormats); err != nil {
		return err
	}

	return nil
}

func (c *_ExpandCommand) run(cmd *cobra.Command, args []string) (err error) {
	var (
		tmpl      utils.TemplateWrapper
		fileCount = 0
	)

	tmpl = utils.NewTemplateWrapper(c.outputFormat)

	if c.templateText != "" {
		if err = tmpl.NewTextTemplate(c.templateText); err != nil {
			return
		}
	} else {
		for _, templateFile := range c.templateFiles {
			if !utils.FileOrDirectoryExists(templateFile) {
				return fmt.Errorf("Template file '%s' does not exist", templateFile)
			}
		}

		// Parse templates
		for _, templateFile := range c.templateFiles {
			if utils.IsDirectory(templateFile) {
				for _, ext := range c.extensions {
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
					if err = tmpl.ParseGlob(pattern); err != nil {
						return
					}
				}
			} else {
				if err = tmpl.ParseFiles(templateFile); err != nil {
					return
				}
				fileCount++
			}
		}
	}

	// If no files matched, then simply return an error
	if fileCount == 0 && c.templateText == "" {
		return fmt.Errorf("No matching file(s)")
	}
	return tmpl.Execute(cmd.OutOrStdout(), c.globalOpts.YamlFile().Data())
}
