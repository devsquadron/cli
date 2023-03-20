/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"
	"net/mail"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/devsquadron/cli/message"
	"github.com/devsquadron/cli/system"
	"github.com/devsquadron/models"
	"github.com/spf13/cobra"
)

var (
	noCacheFlag         bool
	createDeveloperFlag bool
	loginDeveloperFlag  bool
	createTeamFlag      bool
	joinTeamFlag        bool
	cfgPrompts          []system.ConfigPrompt
)

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
				switch ti.Prompt {
				case "Email: ":
					addr, err := mail.ParseAddress(val)
					if err != nil {
						return err
					}
					val = addr.Address
				case "Username: ":
					Cfg.SetToken("")
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
	} else if loginDeveloperFlag {
		err = runDeveloperLogin()
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

	if joinTeamFlag {
		err = TeamClient.JoinTeam(Cfg.Token(), Cfg.Team())
		if err != nil {
			return err
		}
		fmt.Println(message.Green("SUCCESS", fmt.Sprintf("requested to join team %s", Cfg.Team())))
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
		&loginDeveloperFlag, "login-developer", "", false,
		"create login the developer username",
	)
	initCmd.Flags().BoolVarP(
		&createTeamFlag, "create-team", "", false,
		"create a devsquadron team to encapsulate the tasks for your project (after setting local config)",
	)

	initCmd.Flags().BoolVarP(
		&joinTeamFlag, "join-team", "", false,
		"join existing devsquadron team",
	)

	rootCmd.AddCommand(initCmd)
}
