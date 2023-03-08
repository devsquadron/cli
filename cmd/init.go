/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"net/mail"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/devsquadron/cli/system"
	"github.com/spf13/cobra"
)

var (
	noCacheFlag         bool
	createDeveloperFlag bool
	createTeamFlag      bool
	cfgPrompts          []system.ConfigPrompt
)

func runInit() error {
	var (
		err          error
		outMod       tea.Model
		m            system.ConfigurationTextInputModel
		toGetPrompts []system.ConfigPrompt
		existingCfg  string
	)

	cfgPrompts = []system.ConfigPrompt{
		{
			Prompt:      "Username",
			CheckFunc:   Cfg.UsernameE,
			DefaultFunc: Sys.GitUsername,
			SetFunc:     Cfg.SetUsername,
		},
		{
			Prompt:      "Email",
			CheckFunc:   Cfg.EmailE,
			DefaultFunc: Sys.GitEmail,
			SetFunc:     Cfg.SetEmail,
		},
		{
			Prompt: "Team",
			CheckFunc: func() (string, error) {
				// when ds init is called we are either creating or updating a team
				return "", nil
			},
			DefaultFunc: Cfg.TeamE,
			SetFunc:     Cfg.SetTeam,
		},
	}
	for _, cp := range cfgPrompts {
		existingCfg, _ = cp.CheckFunc()
		if existingCfg == "" || noCacheFlag {
			toGetPrompts = append(toGetPrompts, cp)
		}
	}
	m = system.InitialConfigurationTextInputModel(toGetPrompts, noCacheFlag)

	outMod, err = tea.NewProgram(m).Run()
	if err != nil {
		return err
	}

	if m, ok := outMod.(system.ConfigurationTextInputModel); ok {
		for i, ti := range m.Inputs {
			val := ti.Value()
			if val == "" {
				toGetPrompts[i].SetFunc(ti.Placeholder)
			} else {
				if ti.Prompt == "Email: " {
					addr, err := mail.ParseAddress(val)
					if err != nil {
						return err
					}
					val = addr.Address
				}
				toGetPrompts[i].SetFunc(val)
			}
		}
	}
	if createDeveloperFlag {
		err = runDeveloperCreate()
		if err != nil {
			return err
		}
	}
	if createTeamFlag {
		err = runTeamCreate()
		if err != nil {
			return err
		}
	}
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a devsquadron repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func init() {
	initCmd.Flags().BoolVarP(
		&noCacheFlag, "no-cache", "n", false,
		"ignores cache when setting config (so you can revisit initial decisions)",
	)
	initCmd.Flags().BoolVarP(
		&createDeveloperFlag, "create-developer", "", false,
		"create a developer account (using username & email) after setting local config",
	)
	initCmd.Flags().BoolVarP(
		&createTeamFlag, "create-team", "", false,
		"create a devsquadron team to encapsulate the tasks for your project (after setting local config)",
	)

	rootCmd.AddCommand(initCmd)
}
