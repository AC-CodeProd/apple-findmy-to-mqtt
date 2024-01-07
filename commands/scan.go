package commands

import (
	"apple-findmy-to-mqtt/commands/cli"
	"apple-findmy-to-mqtt/core/interfaces"
	"apple-findmy-to-mqtt/infrastructure/config"
	"apple-findmy-to-mqtt/infrastructure/logging"
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
}

func (sC *ScanCommand) GetFlags() interface{} {
	return &ScanCommandWrapper{Path: sC.envPath}
}

func (sC *ScanCommand) Run() cli.ICommandRunner {
	const names = "__scan.go__: Run"
	return func(
		cacheSyncMQTTController interfaces.ICacheSyncMQTTController,
		config config.Config,
		logger logging.Logger,
	) {
		loc, _ := time.LoadLocation(config.TZ)
		time.Local = loc
		logger.Info(fmt.Sprintf("%s | %s", names, "Starting the scan ..."))
		ticker := time.NewTicker(time.Duration(config.ScanTimer) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			logger.Info(fmt.Sprintf("%s | %s", names, "Running scan"))
			cacheSyncMQTTController.Process(config.ForceSync)
		}
	}
}
