/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/devsquadron/cli/message"
	"github.com/devsquadron/models"
	"github.com/spf13/cobra"
)

const (
	ASSIGN_CMD_ARGS = "<id> <developer>"
)

var assignCmd = &cobra.Command{
	Use:   "assign",
	Short: fmt.Sprintf("%s designate the dev to work on the task", ASSIGN_CMD_ARGS),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err       error
			taskIdStr string
			taskId    uint64
			validDev  = false
			toAssign  string
			origTsk   *models.Task
		)

		taskIdStr, toAssign = Sys.GetTwoArgs(args, ASSIGN_CMD_ARGS)
		taskId, err = strconv.ParseUint(taskIdStr, 10, 64)
		if err != nil {
			return err
		}

		if toAssign != "" {
			for _, dv := range Cfg.Developers() {
				if dv == toAssign {
					validDev = true
				}

			}
			if !validDev {
				return errors.New(fmt.Sprintf("%s is not in %s, run 'ds info' to see all developers on the team", toAssign, Cfg.Developers()))
			}
		}

		tm := Cfg.Team()
		origTsk, err = TaskClient.GetTaskById(Cfg.Token(), taskId, tm)
		if err != nil {
			return err
		}
		origTsk.Developer = toAssign

		err = TaskClient.UpdateTask(Cfg.Token(), origTsk, tm)
		if err != nil {
			return err
		}

		message.TaskAbb(origTsk, true)

		return nil
	},

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
		return Cfg.Developers(), cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(assignCmd)
}
