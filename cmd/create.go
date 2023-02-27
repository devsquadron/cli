/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/devsquadron/cli/message"
	"github.com/devsquadron/models"
	"github.com/spf13/cobra"
)

const (
	CREATE_FLAG_NAME_TAG = "tag"
)

var (
	createTagsFlag []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "<title> makes a new task with the specified title",
	Long:  `<title> makes a new task with the specified title`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err     error
			tskBody string
			tsk     = models.Task{}
		)
		if len(args) == 0 {
			return errors.New("expected argument <title>")
		}
		tsk.Title = strings.Join(args, " ")

		if len(createTagsFlag) > 0 {
			tsk.Tags = createTagsFlag
		} else {
			tsk.Tags = append(createTagsFlag, "na")
		}

		// TODO: entry point to integrate template

		tskBody, err = Sys.EditTempMarkdownFileDS("")
		if err != nil {
			return err
		}
		tsk.Criterion = tskBody

		tm := Cfg.Team()
		tsk.ID, err = TaskClient.CreateNewTask(Cfg.Token(), &tsk, tm)
		if err != nil {
			return err
		}
		fmt.Println(message.Green("Success", fmt.Sprintf("%d", tsk.ID)))

		return nil
	},
}

func init() {
	createCmd.Flags().StringArrayVarP(&createTagsFlag, CREATE_FLAG_NAME_TAG, "t", []string{}, "tag to add when creating a task. This flag can be used multiple times.")
	createCmd.RegisterFlagCompletionFunc(CREATE_FLAG_NAME_TAG, func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return Cfg.Tags(), cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(createCmd)
}
