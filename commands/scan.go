package commands

import (
	"apple-findmy-to-mqtt/infrastructure/config"
	"apple-findmy-to-mqtt/infrastructure/logging"
	"apple-findmy-to-mqtt/interfaces/cli"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

type ScanCommand struct {
	envPath string
}

type ScanCommandWrapper struct {
	Path string
}

// create a new run command
func NewScanCommand() *ScanCommand {
	return &ScanCommand{}
}

func (sC *ScanCommand) Short() string {
	return "scan"
}

func (sC *ScanCommand) Setup(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&sC.envPath, "env", "e", "", "Specify the .env file(s).")
	_ = cmd.MarkFlagRequired("env")
}

func (sC *ScanCommand) GetFlags() interface{} {
	return &ScanCommandWrapper{Path: sC.envPath}
}

func (sC *ScanCommand) Run() cli.CommandRunner {
	const names = "__scan.go__: Run"
	return func(
		config config.Config,
		logger logging.Logger,
	) {

		logger.Info(fmt.Sprintf("%s | %s", names, "Starting the scan ..."))
		ticker := time.NewTicker(time.Duration(config.Timer) * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			logger.Info(fmt.Sprintf("%s | %s", names, "Running scan"))
		}
	}
}
