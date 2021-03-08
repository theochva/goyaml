package commands

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

func setExamplesAtEndOfHelp(cmd *cobra.Command) {
	if cmd != nil {
		cmd.SetUsageTemplate(_UsageTemplate)
	}
}

func setUsageReturnCode(cmd *cobra.Command, rc int) {
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

func replaceProgName(format string, a ...interface{}) string {
	str := fmt.Sprintf(format, a...)
	_, progName := filepath.Split(os.Args[0])
	return strings.ReplaceAll(str, "$PROG_NAME", progName)
}
