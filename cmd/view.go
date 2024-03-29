/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/devsquadron/cli/message"

	"github.com/devsquadron/models"

	"github.com/spf13/cobra"
)

var (
	viewFlagRaw bool
)

const (
	VIEW_CMD_ARGS = "<id>"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: fmt.Sprintf("%s view the task info by id", VIEW_CMD_ARGS),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err     error
			taskId  uint64
			origTsk *models.Task
		)

		if len(args) > 0 {
			taskId, err = Sys.GetTaskId(Sys.GetArg(args, VIEW_CMD_ARGS))
			if err != nil {
				return err
			}
		} else {
			taskId, err = Cfg.TaskNumber()
			if err != nil {
				return err
			}
		}

		tm := Cfg.Team()
		origTsk, err = TaskClient.GetTaskById(Cfg.Token(), taskId, tm)
		if err != nil {
			return err
		}

		if viewFlagRaw {
			message.Task(origTsk)
		} else {
			return message.PrettyTask(origTsk)
		}
		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(
		&viewFlagRaw, "raw", "r", false,
		"list tasks assigned to the specified developer",
	)

	rootCmd.AddCommand(viewCmd)
}
