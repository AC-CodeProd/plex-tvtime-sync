package commands

import (
	"encoding/json"
	"fmt"
	"plex-tvtime-sync/interfaces"
	"plex-tvtime-sync/pkg/lib"

	"github.com/spf13/cobra"
)

type CmdCommand struct {
	configPath []string
	jsonString string
}

// create a new run command
func NewCmdCommand() *CmdCommand {
	return &CmdCommand{}
}

func (rC *CmdCommand) Short() string {
	return "cmd"
}

func (rC *CmdCommand) Setup(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&rC.configPath, "config", "c", []string{}, "Specify the configuration file(s).")
	cmd.Flags().StringVarP(&rC.jsonString, "add", "a", "{}", "")
	_ = cmd.MarkFlagRequired("config")
}

func (rC *CmdCommand) GetFlags() []string {
	return rC.configPath
}

func (rC *CmdCommand) Run() lib.CommandRunner {
	const names = "__run.go__: CMD"
	jsonString := rC.jsonString
	return func(
		cmd interfaces.CMD,
		logger lib.Logger,
	) {
		// syncHandler.Run()

		fmt.Println(jsonString)
		type IDs struct {
			PlexID   int64 `json:"plexId"`
			TVTimeID int64 `json:"tvtimeId"`
		}
		var ids IDs
		if err := json.Unmarshal([]byte(jsonString), &ids); err != nil {
			panic(err)
		}
		if err := cmd.AddSpecificPair(&ids.PlexID, &ids.TVTimeID); err != nil {
			panic(err)
		}
		logger.Info(fmt.Sprintf("%s | %s", names, "CMD ..."))
	}
}
