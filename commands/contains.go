package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type _ContainsCmd struct{}

func init() {
	globalOpts.addCommand(
		(&_ContainsCmd{}).createCLICommand(),
		false, // Dont care for yaml validation errors
		false) // Dont skip parsing yaml
}

func (o *_ContainsCmd) createCLICommand() *cobra.Command {
	var cmd = &cobra.Command{
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
		RunE:       o.run,
		Example: replaceProgName(`  $PROG_NAME -f /tmp/foo.yaml contains first.second.third
  $PROG_NAME -f /tmp/foo.yaml has first.second.third
  $PROG_NAME -f /tmp/foo.yaml c first.second.third

  cat /tmp/foo.yaml | $PROG_NAME contains first.second.third
  cat /tmp/foo.yaml | $PROG_NAME has first.second.third
  cat /tmp/foo.yaml | $PROG_NAME c first.second.third`),
	}

	return cmd
}

func (o *_ContainsCmd) run(cmd *cobra.Command, args []string) (err error) {
	var (
		key       = args[0]
		contained bool
	)

	if contained, err = globalOpts.yamlFile.Contains(key); err != nil {
		return
	}

	fmt.Println(contained)
	return
}
