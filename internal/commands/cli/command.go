package cli

import (
	"github.com/spf13/cobra"
)

// AppCommand - the application command interface
type AppCommand interface {
	// GetCliCommand - get the cobra command
	GetCliCommand() *cobra.Command
}

// AppCommandBase - base app command
type AppCommandBase struct {
	cmd *cobra.Command
}

// NewAppCommandBase - create new base app command
func NewAppCommandBase(cmd *cobra.Command) AppCommand {
	return &AppCommandBase{
		cmd: cmd,
	}
}

// GetCliCommand - get the cobra command
func (c *AppCommandBase) GetCliCommand() *cobra.Command { return c.cmd }

// AppSubCommand - sub command interface
type AppSubCommand interface {
	AppCommand
}

// AppSubCommandBase - base app command
type AppSubCommandBase struct {
	AppCommand
}

// NewAppSubCommandBase - create new base app command
func NewAppSubCommandBase(cmd *cobra.Command) AppCommand {
	return &AppSubCommandBase{
		AppCommand: NewAppCommandBase(cmd),
	}
}

// AppRootCommand - root App command interface
type AppRootCommand interface {
	AppCommand
	AddSubCommands(appSubCmds ...AppSubCommand)
}

// AppRootCommandBase - base root command
type AppRootCommandBase struct {
	AppCommand
}

// NewAppRootCommandBase - create new base root command
func NewAppRootCommandBase(cmd *cobra.Command) AppRootCommand {
	return &AppRootCommandBase{
		AppCommand: NewAppCommandBase(cmd),
	}
}

// AddSubCommands - add subcommands
func (c *AppRootCommandBase) AddSubCommands(appSubCmds ...AppSubCommand) {
	for _, appSubCmd := range appSubCmds {
		c.GetCliCommand().AddCommand(appSubCmd.GetCliCommand())
	}
}
