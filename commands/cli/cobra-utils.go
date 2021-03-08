package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const _UsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] [<flags>]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" or "{{.CommandPath}} help [command]" for more information about a command.{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}

`

const _VersionTemplate = `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}{{if .Annotations}}
{{with .Annotations}}{{.Author}}{{end}}{{end}}
`

// SetExamplesAtEndOfUsage - set the examples at the end of the usage text
func SetExamplesAtEndOfUsage(cmd *cobra.Command) {
	if cmd != nil {
		cmd.SetUsageTemplate(_UsageTemplate)
	}
}

// SetUsageReturnCode - the RC to return when usage is displayed
func SetUsageReturnCode(cmd *cobra.Command, rc int) {
	if cmd != nil {
		// Make the help function return a non-zero rc
		helpFunc := cmd.HelpFunc()
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			// call the default help func
			helpFunc(cmd, args)
			// Exit with non-zero
			os.Exit(rc)
		})
	}
}

// SetVersionWithAuthor - show author info in the version
func SetVersionWithAuthor(cmd *cobra.Command, author string) {
	if author != "" {
		if cmd.Annotations == nil {
			cmd.Annotations = map[string]string{}
		}
		cmd.Annotations["Author"] = author
		// versionTemplate := cmd.VersionTemplate()
		// cmd.SetVersionTemplate(_VersionTemplate)
	}
	cmd.SetVersionTemplate(_VersionTemplate)
}

// ReplaceProgName - replace the program name ("$PROG_NAME") in a string
func ReplaceProgName(format string, a ...interface{}) string {
	str := fmt.Sprintf(format, a...)
	_, progName := filepath.Split(os.Args[0])
	return strings.ReplaceAll(str, "$PROG_NAME", progName)
}
