package cli

// App - application wrapper for root command
type App struct {
	rootCommand AppRootCommand
}

// NewApp - create new app
func NewApp(rootCommand AppRootCommand) *App {
	return &App{
		rootCommand: rootCommand,
	}
}

// GetRootCommand - get the root command
func (a *App) GetRootCommand() AppRootCommand { return a.rootCommand }

// Execute - run the app
func (a *App) Execute() error {
	return a.rootCommand.GetCliCommand().Execute()
	// if err := a.rootCommand.GetCliCommand().Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}

// if err := globalOpts.rootCmd.Execute(); err != nil {
// 	if _, isValidationErr := err.(validationError); !isValidationErr {
// 		fmt.Println(err)
// 	}
// 	os.Exit(1)
// }
