package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type _DeleteCmd struct{}

func init() {
	globalOpts.addCommand(
		(&_DeleteCmd{}).createCLICommand(),
		false, // Dont care for yaml validation errors
		false) // Dont skip parsing yaml
}

func (o *_DeleteCmd) createCLICommand() *cobra.Command {
	var cmd = &cobra.Command{
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
		RunE:       o.run,
		Example: replaceProgName(`  $PROG_NAME -f /tmp/foo.yaml delete first.second.third
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

	return cmd
}

func (o *_DeleteCmd) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key      = args[0]
		yamlText string
		deleted  bool
	)

	if deleted, err = globalOpts.yamlFile.Delete(key); err != nil {
		return
	}

	if deleted {
		// If YAML read from stdin, then "Save" will output result
		if err = globalOpts.yamlFile.Save(); err != nil {
			return err
		}
	}

	// If YAML not read from stdin, then print the delete result
	if !globalOpts.pipe {
		fmt.Println(deleted)
	} else {
		// Else, YAML read from stdin. If not deleted,
		// then nothing printed, so dump the YAML
		if !deleted {
			if yamlText, err = globalOpts.yamlFile.Text(); err != nil {
				return err
			}
			fmt.Println(yamlText)
		}
	}
	return
}
