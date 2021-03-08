package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type _ValidateCmd struct {
	details bool
}

func init() {
	globalOpts.addCommand(
		(&_ValidateCmd{}).createCLICommand(),
		true,  // Do care for yaml validation errors
		false) // Dont skip parsing yaml
}

func (o *_ValidateCmd) createCLICommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "validate [-d|--details]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"v"},
		Short:                 "Validate the yaml syntax",
		Long:                  "Validate the  yaml syntax. It either outputs 'true', 'false' or the validation msg.",
		Args:                  cobra.NoArgs,
		RunE:                  o.run,
		Example: replaceProgName(`  $PROG_NAME -f /tmp/foo.yaml validate
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

	cmd.Flags().BoolVarP(
		&o.details,
		_flagDetails, _flagDetailsShort, false,
		"Prints the parsing error instead of 'false'.  If valid, it outputs nothing",
	)

	return cmd
}

func (o *_ValidateCmd) run(cmd *cobra.Command, args []string) (err error) {
	valid := (globalOpts.yamlValidationErr == nil)

	if !valid && o.details {
		fmt.Println(globalOpts.yamlValidationErr.Error())
		return
	}
	fmt.Println(valid)
	return
}
