package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theochva/goyaml/internal/commands/cli"
)

type _ContainsCommand struct {
	cli.AppSubCommand

	globalOpts GlobalOptions
}

func newContainsCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_ContainsCommand{
		globalOpts: globalOpts,
	}

	cliCmd := &cobra.Command{
		Use:                   "contains <key>",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"c", "has"},
		Short:                 "Check if a value is contained in the yaml",
		Long:                  "Check if a value is contained in the yaml. It simply outputs 'true' or 'false'.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires the 'key' to check")
			}
			return nil
		},
		ArgAliases: []string{"key"},
		RunE:       subCmd.run,
		Example: cli.ReplaceProgName(`  $PROG_NAME -f /tmp/foo.yaml contains first.second.third
  $PROG_NAME -f /tmp/foo.yaml has first.second.third
  $PROG_NAME -f /tmp/foo.yaml c first.second.third

  cat /tmp/foo.yaml | $PROG_NAME contains first.second.third
  cat /tmp/foo.yaml | $PROG_NAME has first.second.third
  cat /tmp/foo.yaml | $PROG_NAME c first.second.third`),
	}

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)
	return subCmd
}

func (c *_ContainsCommand) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key       = args[0]
		contained bool
	)

	if contained, err = c.globalOpts.YamlFile().Contains(key); err != nil {
		return
	}

	cmd.Println(contained)
	return
}
