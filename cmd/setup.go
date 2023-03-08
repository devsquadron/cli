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

var developerLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "login a developer",
	Long:  `login a developer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err error
			dev models.Developer
			tkn string
		)

		dev.Name = Cfg.Username()

		dev.Password, err = Sys.GetPassword()
		if err != nil {
			return err
		}

		tkn, err = DeveloperClient.LoginDeveloper(&dev)
		if err != nil {
			return err
		}

		Cfg.SetToken(tkn)
		fmt.Println()
		fmt.Println(message.Green("SUCCESS", "Saved token to config."))

		return nil
	},
}

func runDeveloperCreate() error {
	var (
		err error
		dev models.Developer
		tkn string
	)

	dev.Name = Cfg.Username()
	dev.Email = Cfg.Email()

	dev.Password, dev.ConfirmPassword, err = Sys.GetConfirmedPassword()
	if err != nil {
		return err
	}

	tkn, err = DeveloperClient.CreateNewDeveloper(&dev)
	if err != nil {
		return err
	}

	Cfg.SetToken(tkn)
	fmt.Println()
	fmt.Println(message.Green("SUCCESS", "Saved token to config."))

	return nil
}

var developerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a developer account using details from init command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeveloperCreate()
	},
}

var developerCmd = &cobra.Command{
	Use:   "developer",
	Short: "setting up developers",
	Long:  `setting up developers`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var teamGrowCmd = &cobra.Command{
	Use:   "grow",
	Short: "<name> add a developer to a team by username",
	Long:  `<name> add a developer to a team by username`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("expected one argument <name>")
		}
		var err error
		err = TeamClient.GrowTeam(
			Cfg.Token(),
			Cfg.Team(),
			&models.Developer{Name: args[0]},
		)
		if err != nil {
			return err
		}
		fmt.Println(message.Green("[SUCCESS]", fmt.Sprintf("added %s to team %s", args[0], Cfg.Team())))
		return nil
	},
}

func runTeamCreate() error {
	var (
		tmNm string
		err  error
	)
	tmNm = Cfg.Team()
	tm := models.Team{Name: tmNm}
	err = TeamClient.CreateNewTeam(
		&tm,
		Cfg.Token(),
	)
	if err != nil {
		return err
	}
	fmt.Println(message.Green("SUCCESS", fmt.Sprintf("created new team '%s'", tm.Name)))
	return nil
}

var teamCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new team using the name provided with ds init",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTeamCreate()
	},
}

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "setting up teams",
	Long:  `setting up teams`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "actions related to setting up",
	Long:  `actions related to setting up`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	developerCmd.AddCommand(developerCreateCmd)
	developerCmd.AddCommand(developerLoginCmd)
	setupCmd.AddCommand(developerCmd)

	teamCmd.AddCommand(teamCreateCmd)
	teamCmd.AddCommand(teamGrowCmd)
	setupCmd.AddCommand(teamCmd)

	rootCmd.AddCommand(setupCmd)
}
