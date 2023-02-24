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
	EMPTY_TAGS_NA   = "na"
	TAG_CMD_ARGS    = "<id> <tag>"
	TAG_FLAG_DELETE = "delete"
)

var (
	tagDeleteFlag bool
)

// TODO the methods filterTag and addTag should be moved to backend
func filterTag(tags []string, toRmTg string) []string {
	i := 0
	for _, ctg := range tags {
		if ctg != toRmTg {
			tags[i] = ctg
			i++
		}
	}
	tags = tags[:i]
	if len(tags) == 0 {
		tags = append(tags, EMPTY_TAGS_NA)
	}
	return tags
}

func addTag(tags []string, toAdTg string) []string {
	noNa := true
	for i, ctg := range tags {
		if ctg == EMPTY_TAGS_NA {
			tags[i] = toAdTg
			noNa = false // found na
			break
		}
	}
	if noNa {
		tags = append(tags, toAdTg)
	}
	return tags
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: fmt.Sprintf("%s give the task a tag", TAG_CMD_ARGS),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err        error
			taskId     uint64
			tag, idStr string
			origTsk    *models.Task
		)

		idStr, tag = Sys.GetTwoArgs(args, TAG_CMD_ARGS)
		taskId, err = Sys.GetTaskId(idStr)
		if err != nil {
			return err
		}

		tm := Cfg.Team()
		origTsk, err = TaskClient.GetTaskById(Cfg.Token(), taskId, tm)
		if err != nil {
			return err
		}

		// TODO this validation could be moved inside system
		if tag == "" {
			return errors.New("Cannot have an empty <tag> argument")
		}
		if tagDeleteFlag {
			origTsk.Tags = filterTag(origTsk.Tags, tag)
		} else {
			origTsk.Tags = addTag(origTsk.Tags, tag)
		}
		err = TaskClient.UpdateTask(Cfg.Token(), origTsk, tm)
		if err != nil {
			return err
		}

		message.Task(origTsk)
		return nil
	},
	// TODO test this out
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
		return Cfg.Tags(), cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	tagCmd.Flags().BoolVarP(&tagDeleteFlag, TAG_FLAG_DELETE, "d", false, "remove a tag from a task")

	rootCmd.AddCommand(tagCmd)
}
