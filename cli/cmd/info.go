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

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "get information about your user and team",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err    error
			tmInfo *models.Team
		)

		fmt.Println(message.Green("Username", Cfg.Username()))
		fmt.Println("[Email]", Cfg.Email())
		tm := Cfg.Team()
		tmInfo, err = TeamClient.InfoTeam(Cfg.Token(), tm)
		if err != nil {
			return err
		}

		fmt.Println(message.Green("Team", tm))
		fmt.Println("[Developers]", tmInfo.Developers)
		fmt.Println("[Task Count]", tmInfo.TaskCount)
		Cfg.SetDevelopers(tmInfo.Developers)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
