/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/devsquadron/cli/message"

	"github.com/devsquadron/models"

	"github.com/spf13/cobra"
)

const (
	EDIT_CMD_ARGS         = "<id>"
	EDIT_FLAG_NAME_STATUS = "status"
	EDIT_FLAG_NAME_TITLE  = "title"
)

var (
	editStatusFlag string
	editTitleFlag  string
)

func updateTaskBody(origTsk *models.Task) error {
	var (
		err          error
		newCriterion string
	)
	newCriterion, err = Sys.EditTempMarkdownFileDS(origTsk.Criterion)
	if err != nil {
		return err
	}
	if newCriterion == origTsk.Criterion {
		return errors.New("there were no updates applied to the task body.")
	}
	origTsk.Criterion = newCriterion
	return nil
}

func processTheFlags(tsk *models.Task) error {
	if editStatusFlag != "" {
		if models.IsStatus(editStatusFlag) {
			tsk.Status = editStatusFlag
		} else {
			return errors.New(fmt.Sprintf("The status %s is not a valid status", editStatusFlag))
		}
	}
	if editTitleFlag != "" {
		// TODO: check here if title is a valid git branch
		tsk.Title = editTitleFlag
	}
	return nil
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: fmt.Sprintf("%s edit a task", EDIT_CMD_ARGS),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err     error
			origTsk *models.Task
			taskId  uint64
		)
		taskId, err = Sys.GetTaskId(Sys.GetArg(args, EDIT_CMD_ARGS))
		if err != nil {
			return err
		}

		tm := Cfg.Team()
		origTsk, err = TaskClient.GetTaskById(Cfg.Token(), taskId, tm)
		if err != nil {
			return err
		}

		if editStatusFlag != "" || editTitleFlag != "" {
			err = processTheFlags(origTsk)
			if err != nil {
				return err
			}
		} else {
			err = updateTaskBody(origTsk)
			if err != nil {
				return err
			}
		}

		err = TaskClient.UpdateTask(Cfg.Token(), origTsk, tm)
		if err != nil {
			return err
		}
		fmt.Println(message.Green("Success", "task updated successfully."))

		return nil
	},
}

func init() {
	editCmd.Flags().StringVarP(&editStatusFlag, EDIT_FLAG_NAME_STATUS, "s", "", "changes the status for a task")
	editCmd.RegisterFlagCompletionFunc(EDIT_FLAG_NAME_STATUS, func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return models.Statuses, cobra.ShellCompDirectiveNoFileComp
	})

	editCmd.Flags().StringVarP(&editTitleFlag, EDIT_FLAG_NAME_TITLE, "t", "", "changes the title for a task")

	rootCmd.AddCommand(editCmd)
}
