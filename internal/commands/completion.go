package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/theochva/goyaml/internal/commands/cli"
)

type _CompletionCommand struct {
	cli.AppSubCommand
}

func init() {
	registerCommand(func(_ GlobalOptions) cli.AppSubCommand {
		return _CompletionCommand{
			AppSubCommand: cli.NewAppSubCommandBase(&cobra.Command{
				Use:                   "completion [bash|zsh|fish|powershell]",
				DisableFlagsInUseLine: true,
				Hidden:                true,
				Annotations: map[string]string{
					_CmdOptValidationAware: _CmdOptValueTrue,
					_CmdOptSkipParsing:     _CmdOptValueTrue,
				},
				Short:     "Generate completion script",
				ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
				Args:      cobra.ExactValidArgs(1),
				Run: func(cmd *cobra.Command, args []string) {
					switch args[0] {
					case "bash":
						cmd.Root().GenBashCompletion(os.Stdout)
					case "zsh":
						cmd.Root().GenZshCompletion(os.Stdout)
					case "fish":
						cmd.Root().GenFishCompletion(os.Stdout, true)
					case "powershell":
						cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
					}
				},
				Long: cli.ReplaceProgName(`To load completions:

Bash:
				
  $ source <($PROG_NAME completion bash)
				
  # To load completions for each session, execute once:
  # Linux:
  $ $PROG_NAME completion bash > /etc/bash_completion.d/$PROG_NAME
  # macOS:
  $ $PROG_NAME completion bash > /usr/local/etc/bash_completion.d/$PROG_NAME
				
Zsh:
				
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
				
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
				
  # To load completions for each session, execute once:
  $ $PROG_NAME completion zsh > "${fpath[1]}/_$PROG_NAME"
				
  # You will need to start a new shell for this setup to take effect.
				
fish:
				
  $ $PROG_NAME completion fish | source
				
  # To load completions for each session, execute once:
  $ $PROG_NAME completion fish > ~/.config/fish/completions/$PROG_NAME.fish
				
PowerShell:
				
  PS> $PROG_NAME completion powershell | Out-String | Invoke-Expression
				
  # To load completions for every new session, run:
  PS> $PROG_NAME completion powershell > $PROG_NAME.ps1
  # and source this file from your PowerShell profile.
`),
			}),
		}
	})
}
