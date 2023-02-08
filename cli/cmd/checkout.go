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
	CHECKOUT_CMD_ARGS = "<id>"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: fmt.Sprintf("%s checkout a git branch, assign to yourself, move to %s", CHECKOUT_CMD_ARGS, models.Status.Developing),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err     error
			taskId  uint64
			dsUsrNm string
			origTsk *models.Task
		)

		taskId, err = Sys.GetTaskId(Sys.GetArg(args, CHECKOUT_CMD_ARGS))
		if err != nil {
			return err
		}

		dsUsrNm, err = Cfg.UsernameE()
		if err != nil {
			return err
		}

		tm := Cfg.Team()
		origTsk, err = TaskClient.GetTaskById(Cfg.Token(), taskId, tm)
		if err != nil {
			return err
		}
		origTsk.Developer = dsUsrNm
		origTsk.Status = models.Status.Developing

		err = Sys.CheckoutBranch(origTsk)
		if err != nil {
			fmt.Println(
				message.Yellow("Warning", "did not checkout git branch"),
			)
			fmt.Println(err)
		}

		err = TaskClient.UpdateTask(Cfg.Token(), origTsk, tm)
		if err != nil {
			return err
		}

		message.Task(origTsk)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
