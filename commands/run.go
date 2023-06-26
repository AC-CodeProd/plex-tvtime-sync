package commands

import (
	"fmt"
	"plex-tvtime-sync/domain/process"
	"plex-tvtime-sync/pkg/lib"

	"github.com/spf13/cobra"
)

type RunCommand struct {
	configPath []string
}

// create a new run command
func NewRunCommand() *RunCommand {
	return &RunCommand{}
}

func (rC *RunCommand) Short() string {
	return "run"
}

func (rC *RunCommand) Setup(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&rC.configPath, "config", "c", []string{}, "Specify the configuration file(s).")
	_ = cmd.MarkFlagRequired("config")
}

func (rC *RunCommand) GetFlags() []string {
	return rC.configPath
}

func (rC *RunCommand) Run() lib.CommandRunner {
	const names = "__run.go__: Run"
	return func(
		syncHandler process.SyncProcess,
		logger lib.Logger,
	) {
		syncHandler.Run()
		logger.Info(fmt.Sprintf("%s | %s", names, "Starting the bot ..."))
	}
}
