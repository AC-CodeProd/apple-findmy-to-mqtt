package commands

import (
	"apple-findmy-to-mqtt/infrastructure/config"
	"apple-findmy-to-mqtt/infrastructure/logging"
	"apple-findmy-to-mqtt/interfaces/cli"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var cmds = map[string]cli.Command{
	"scan": NewScanCommand(),
}

// get a list of sub commands
func GetSubCommands(opt fx.Option) []*cobra.Command {
	var subCmds []*cobra.Command

	for name, cmd := range cmds {
		subCmds = append(subCmds, wrapSubCommand(name, cmd, opt))
	}

	return subCmds
}

func wrapSubCommand(name string, cmd cli.Command, opt fx.Option) *cobra.Command {
	const names = "__commands.go__ : wrapSubCommand"
	subCmd := &cobra.Command{
		Use:   name,
		Short: cmd.Short(),
		Run: func(c *cobra.Command, args []string) {
			serverCommandwrapper, _ := cmd.GetFlags().(*ScanCommandWrapper)
			if err := config.SetupConfig(serverCommandwrapper.Path); err != nil {
				panic(fmt.Sprintf("%s | %s", names, err))
			}
			logger := logging.GetLogger()
			opts := fx.Options(
				fx.WithLogger(func() fxevent.Logger {
					return logger.GetFxLogger()
				}),
				fx.Invoke(cmd.Run()),
			)
			ctx := context.Background()
			app := fx.New(opt, opts)
			if err := app.Start(ctx); err != nil {
				logger.Fatal(fmt.Sprintf("%s | %s", names, err))
				panic(fmt.Sprintf("%s | %s", names, err))
			}
			if err := app.Stop(ctx); err != nil {
				logger.Fatal(fmt.Sprintf("%s | %s", names, err))
				panic(fmt.Sprintf("%s | %s", names, err))
			}
		},
	}

	cmd.Setup(subCmd)
	return subCmd
}
