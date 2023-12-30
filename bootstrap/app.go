package bootstrap

import (
	"apple-findmy-to-mqtt/commands"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "",
	Short:            "",
	Long:             "",
	TraverseChildren: true,
}

type App struct {
	*cobra.Command
}

func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(commands.GetSubCommands(CommonModules)...)

	return cmd
}

var RootApp = NewApp()
