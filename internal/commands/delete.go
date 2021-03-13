package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theochva/goyaml/internal/commands/cli"
)

type _DeleteCommand struct {
	cli.AppSubCommand

	globalOpts GlobalOptions
}

func newDeleteCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_DeleteCommand{
		globalOpts: globalOpts,
	}

	cliCmd := &cobra.Command{
		Use:                   "delete <key>",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"d", "del", "remove", "rm"},
		Short:                 "Delete a value from the yaml",
		Long: `Delete a value from the yaml. If reading from stdin, it outputs the updated YAML. If reading
from a file, it simply outputs 'true' or 'false' to indicate whether the value was deleted.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires the 'key' to delete")
			}
			return nil
		},
		ArgAliases: []string{"key"},
		RunE:       subCmd.run,
		Example: cli.ReplaceProgName(`  $PROG_NAME -f /tmp/foo.yaml delete first.second.third
  $PROG_NAME -f /tmp/foo.yaml del first.second.third
  $PROG_NAME -f /tmp/foo.yaml d first.second.third
  $PROG_NAME -f /tmp/foo.yaml remove first.second.third
  $PROG_NAME -f /tmp/foo.yaml rm first.second.third

  When piping YAML text, the updated YAML is printed to stdout:
    cat /tmp/foo.yaml | $PROG_NAME delete first.second.third
    cat /tmp/foo.yaml | $PROG_NAME del first.second.third
    cat /tmp/foo.yaml | $PROG_NAME d first.second.third
    cat /tmp/foo.yaml | $PROG_NAME remove first.second.third
    cat /tmp/foo.yaml | $PROG_NAME rm first.second.third`),
	}

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)
	return subCmd
}

func (c *_DeleteCommand) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key      = args[0]
		yamlText string
		deleted  bool
	)

	if deleted, err = c.globalOpts.YamlFile().Delete(key); err != nil {
		return
	}

	if deleted {
		// If YAML read from stdin, then "Save" will output result
		if err = c.globalOpts.YamlFile().Save(); err != nil {
			return err
		}
	}

	// If YAML not read from stdin, then print the delete result
	if !c.globalOpts.IsPipe() {
		cmd.Println(deleted)
	} else {
		// Else, YAML read from stdin. If not deleted,
		// then nothing printed, so dump the YAML
		if !deleted {
			if yamlText, err = c.globalOpts.YamlFile().Text(); err != nil {
				return err
			}
			cmd.Println(yamlText)
		}
	}
	return
}
