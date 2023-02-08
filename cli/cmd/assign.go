/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/devsquadron/ds/message"

	"github.com/devsquadron/models"

	"github.com/spf13/cobra"
)

const (
	ASSIGN_CMD_ARGS            = "<id>"
	ASSIGN_FLAG_NAME_DEVELOPER = "developer"
)

var (
	assignDevFlag string
)

var assignCmd = &cobra.Command{
	Use:   "assign",
	Short: fmt.Sprintf("%s designate the dev to work on the task", ASSIGN_CMD_ARGS),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err     error
			taskId  uint64
			dsUsrNm string
			origTsk *models.Task
		)

		taskId, err = Sys.GetTaskId(Sys.GetArg(args, ASSIGN_CMD_ARGS))
		if err != nil {
			return err
		}

		if assignDevFlag != "" {
			dsUsrNm = assignDevFlag
		} else {
			dsUsrNm, err = Cfg.UsernameE()
			if err != nil {
				return err
			}
		}

		tm := Cfg.Team()
		origTsk, err = TaskClient.GetTaskById(Cfg.Token(), taskId, tm)
		if err != nil {
			return err
		}
		origTsk.Developer = dsUsrNm

		err = TaskClient.UpdateTask(Cfg.Token(), origTsk, tm)
		if err != nil {
			return err
		}

		message.Task(origTsk)

		return nil
	},
}

func init() {
	assignCmd.Flags().StringVarP(
		&assignDevFlag, ASSIGN_FLAG_NAME_DEVELOPER, "d", "",
		"assign task to the specified developer",
	)
	assignCmd.RegisterFlagCompletionFunc(ASSIGN_FLAG_NAME_DEVELOPER, func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return Cfg.Developers(), cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(assignCmd)
}
