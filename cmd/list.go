/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/devsquadron/cli/message"

	"github.com/devsquadron/models"

	"github.com/spf13/cobra"
)

const (
	LIST_FLAG_NAME_TAG       = "tag"
	LIST_FLAG_NAME_DEVELOPER = "developer"
	LIST_FLAG_NAME_STATUS    = "status"
)

var (
	listDevFlag    string
	listTagFlag    string
	listStatusFlag []string
)

func adjustTagsForPrinting() {
	if len(listStatusFlag) == 0 {
		listStatusFlag = []string{models.Status.New, models.Status.Developing}
	}
	if listTagFlag == "" {
		listTagFlag = "any"
	}
}

func listAllTasks() error {
	var (
		err     error
		allTsks *[]models.Task
		tsk     models.Task
	)

	if len(listStatusFlag) == 0 && listDevFlag != "" {
		listStatusFlag = []string{models.Status.New, models.Status.Developing, models.Status.Review}
	}
	tm := Cfg.Team()
	allTsks, err = TaskClient.GetTasks(Cfg.Token(), tm, listTagFlag, listStatusFlag, listDevFlag)
	if err != nil {
		return err
	}

	// ----------- PAST THIS POINT IS JUST FOR LISTING
	if listTagFlag != "" || listDevFlag != "" || len(listStatusFlag) > 0 {
		adjustTagsForPrinting()
		message.ListTasksByContext(allTsks, listTagFlag, listDevFlag)
	} else {
		adjustTagsForPrinting()
		message.Header()
		fmt.Println(fmt.Sprintf("status %s | tag %s | dev %s", listStatusFlag, listTagFlag, listDevFlag))
		for _, tsk = range *allTsks {
			message.TaskAbb(&tsk, true)
		}
	}

	return nil
}

func listTagCount() error {
	var (
		err        error
		tgTskDst   *[]models.TagDistribution
		newTgCache = []string{}
	)
	tgTskDst, err = TaskClient.GetTagTaskDistribution(Cfg.Token(), Cfg.Team())
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("status %s", []string{models.Status.New, models.Status.Developing}))
	message.Header()
	for _, tgD := range *tgTskDst {
		fmt.Println(tgD.Count, tgD.Tag)
		newTgCache = append(newTgCache, tgD.Tag)
	}
	Cfg.SetTags(newTgCache)
	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if listDevFlag != "" || listTagFlag != "" || len(listStatusFlag) > 0 {
			return listAllTasks()
		}
		return listTagCount()
	},
}

func init() {
	listCmd.Flags().StringVarP(
		&listDevFlag, LIST_FLAG_NAME_DEVELOPER, "d", "",
		"list tasks assigned to the specified developer",
	)
	listCmd.RegisterFlagCompletionFunc(LIST_FLAG_NAME_DEVELOPER, func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return Cfg.Developers(), cobra.ShellCompDirectiveNoFileComp
	})

	listCmd.Flags().StringVarP(
		&listTagFlag, LIST_FLAG_NAME_TAG, "t", "",
		"list the tasks which have a given tag",
	)
	listCmd.RegisterFlagCompletionFunc(LIST_FLAG_NAME_TAG, func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return Cfg.Tags(), cobra.ShellCompDirectiveNoFileComp
	})

	listCmd.Flags().StringSliceVarP(&listStatusFlag, LIST_FLAG_NAME_STATUS, "s", []string{}, "queries for tasks with the specified status")
	listCmd.RegisterFlagCompletionFunc(LIST_FLAG_NAME_STATUS, func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		var (
			comps, alryCompd, rets, flRets []string
			lastComp, compPrefix           string
		)
		if strings.Contains(toComplete, ",") {
			comps = strings.Split(toComplete, ",")
			alryCompd = comps[0 : len(comps)-1]
			lastComp = comps[len(comps)-1]

		out:
			for _, stts := range models.Statuses {
				for _, alryStts := range alryCompd {
					if stts == alryStts {
						continue out
					}
				}
				if strings.HasPrefix(stts, lastComp) {
					rets = append(rets, stts)
				}
			}
			compPrefix = strings.Join(alryCompd, ",")
			for _, iRet := range rets {
				flRets = append(flRets, compPrefix+","+iRet)
			}
			return flRets, cobra.ShellCompDirectiveNoFileComp
		}
		return models.Statuses, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(listCmd)
}
