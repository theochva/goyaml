package commands

import (
	"github.com/spf13/cobra"
	"github.com/theochva/goyaml/commands/cli"
)

type _ValidateCommand struct {
	cli.AppSubCommand
	ValidationErrorAwareCommand

	globalOpts GlobalOptions
	details    bool
}

// newValidateCommand - create the "validate" subcommand
func newValidateCommand(globalOpts GlobalOptions) cli.AppSubCommand {
	subCmd := &_ValidateCommand{
		globalOpts: globalOpts,
	}

	cliCmd := &cobra.Command{
		Use:                   "validate [-d|--details]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"v"},
		Short:                 "Validate the yaml syntax",
		Long:                  "Validate the  yaml syntax. It either outputs 'true', 'false' or the validation msg.",
		Args:                  cobra.NoArgs,
		RunE:                  subCmd.run,
		Example: cli.ReplaceProgName(`  $PROG_NAME -f /tmp/foo.yaml validate
  $PROG_NAME -f /tmp/foo.yaml validate --details
  $PROG_NAME -f /tmp/foo.yaml validate -d
  $PROG_NAME -f /tmp/foo.yaml v
  $PROG_NAME -f /tmp/foo.yaml v --details
  $PROG_NAME -f /tmp/foo.yaml v -d

  cat /tmp/foo.yaml | $PROG_NAME validate
  cat /tmp/foo.yaml | $PROG_NAME validate --details
  cat /tmp/foo.yaml | $PROG_NAME validate -d
  cat /tmp/foo.yaml | $PROG_NAME v
  cat /tmp/foo.yaml | $PROG_NAME v --details
  cat /tmp/foo.yaml | $PROG_NAME v -d`),
	}

	cliCmd.Flags().BoolVarP(
		&subCmd.details,
		_flagDetails, _flagDetailsShort, false,
		"Prints the parsing error instead of 'false'.  If valid, it outputs nothing",
	)

	subCmd.AppSubCommand = cli.NewAppSubCommandBase(cliCmd)
	return subCmd
}

// IsValidationAware - implementing this method to indicate that this command cares for validation errors
func (c *_ValidateCommand) IsValidationAware() bool { return true }

func (c *_ValidateCommand) run(cmd *cobra.Command, args []string) (err error) {
	valid := (c.globalOpts.ValidationError() == nil)

	if !c.details {
		cmd.Println(valid)
	} else if !valid {
		cmd.Println(c.globalOpts.ValidationError().Error())
	}
	return
}
