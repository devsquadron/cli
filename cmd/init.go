/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"github.com/spf13/cobra"
)

func getOrSetConfig(
	prmpt string,
	checkFun func() (string, error),
	defFun func() (string, error),
	setFun func(string),
) error {
	var (
		err                      error
		existingCfg, gitCfg, inp string
	)
	existingCfg, _ = checkFun()
	if existingCfg == "" {
		gitCfg, _ = defFun()
		inp, err = Sys.GetText(prmpt, gitCfg)
		if err != nil {
			return err
		}
		setFun(inp)
	}
	return nil
}

func runInit() error {
	var (
		err error
	)
	err = getOrSetConfig(
		"Username",
		Cfg.UsernameE,
		Sys.GitUsername,
		Cfg.SetUsername,
	)
	if err != nil {
		return err
	}

	err = getOrSetConfig(
		"Email",
		Cfg.EmailE,
		Sys.GitEmail,
		Cfg.SetEmail,
	)
	if err != nil {
		return err
	}

	return getOrSetConfig(
		"Team",
		func() (string, error) {
			// when ds init is called we are either creating or updating a team
			return "", nil
		},
		Cfg.TeamE,
		Cfg.SetTeam,
	)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a devsquadron repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err error
		)

		err = runInit()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
