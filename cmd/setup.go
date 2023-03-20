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

func runDeveloperLogin() error {
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
}

var developerLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "login a developer",
	Long:  `login a developer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeveloperLogin()
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
	Use:   "create-developer",
	Short: "create a developer account using details from init command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeveloperCreate()
	},
}

func approveDenyValidArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return Cfg.RequestingDevelopers(), cobra.ShellCompDirectiveNoFileComp
	}
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}

func runApproveDenyRequest(args []string, appDy models.ApproveDeny) error {
	var (
		err   error
		dev   string
		chDev string
	)

	if len(args) != 1 {
		return errors.New("must pass a requesting developer, run 'ds info' to see requests")
	}
	dev = args[0]

	err = errors.New(fmt.Sprintf(
		"can only approve / deny requesting developers. %s not in %s. run 'ds info' to see requests",
		dev,
		Cfg.RequestingDevelopers(),
	))
	for _, chDev = range Cfg.RequestingDevelopers() {
		if chDev == dev {
			err = nil
		}
	}
	if err != nil {
		return err
	}

	if appDy == models.APPROVE {
		for _, chDev = range Cfg.Developers() {
			if chDev == dev {
				return errors.New(fmt.Sprintf(
					"cannot add developer %s, already on team %s",
					dev,
					Cfg.Team(),
				))
			}
		}
	}

	err = TeamClient.RespondJoinRequest(
		Cfg.Token(),
		Cfg.Team(),
		&models.RespondJoinRequestReq{RequestDev: dev, ApproveDeny: appDy},
	)
	if err != nil {
		return err
	}

	msg := "denied access"
	if appDy == models.APPROVE {
		msg = "given access"
	}
	fmt.Println(message.Green("SUCCESS", fmt.Sprintf(
		"Developer %s was %s to team %s",
		dev,
		msg,
		Cfg.Team(),
	)))

	return nil
}

var requestsDenyCmd = &cobra.Command{
	Use:               "deny",
	Short:             "<developer> disallow joining the team",
	ValidArgsFunction: approveDenyValidArgsFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApproveDenyRequest(args, models.DENY)
	},
}

var requestsApproveCmd = &cobra.Command{
	Use:               "approve",
	Short:             "<developer> add requesting developer to the team",
	ValidArgsFunction: approveDenyValidArgsFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApproveDenyRequest(args, models.APPROVE)
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "actions related to setting up",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	setupCmd.AddCommand(developerCreateCmd)
	setupCmd.AddCommand(developerLoginCmd)

	setupCmd.AddCommand(requestsApproveCmd)
	setupCmd.AddCommand(requestsDenyCmd)

	rootCmd.AddCommand(setupCmd)
}
