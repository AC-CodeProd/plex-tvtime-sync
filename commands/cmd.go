package commands

import (
	"encoding/json"
	"fmt"
	"plex-tvtime-sync/interfaces"
	"plex-tvtime-sync/pkg/lib"
	"strconv"

	"github.com/spf13/cobra"
)

type CmdCommand struct {
	configPath  []string
	jsonString  string
	storageList bool
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
	cmd.Flags().StringVar(&rC.jsonString, "storage:add", "{}", "Manual addition of plexid and tvtimeId mapping.")
	cmd.Flags().BoolVar(&rC.storageList, "storage:list", true, "List of media already mapped.")
	_ = cmd.MarkFlagRequired("config")
}

func (rC *CmdCommand) GetFlags() *map[string]any {
	flags := make(map[string]any)
	// flags["configPath"] = strings.Join(rC.configPath, "|")
	flags["configPath"] = rC.configPath
	flags["jsonString"] = rC.jsonString
	flags["storageList"] = strconv.FormatBool(rC.storageList)
	return &flags
}

func (rC *CmdCommand) Run() lib.CommandRunner {
	const names = "__run.go__: CMD"
	jsonString := rC.jsonString
	storageList := rC.storageList
	return func(
		cmd interfaces.CMD,
		logger lib.Logger,
	) {

		if storageList && jsonString == "" {
			displaySpecificPair(&cmd)
		} else {
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
			displaySpecificPair(&cmd)
		}
		logger.Debug(fmt.Sprintf("%s | %s", names, "CMD ..."))
	}
}

func displaySpecificPair(cmd *interfaces.CMD) {
	// fmt.Println(cmd.GetAllSpecificPair())
	if data, err := cmd.GetAllSpecificPair(); err == nil {
		for plexId, tvtimeId := range data {
			fmt.Printf("PlexId: %d, TVTimeId: %d\n", plexId, tvtimeId)
		}
	}
}
