/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/devsquadron/requests"

	"github.com/devsquadron/cli/configuration"
	"github.com/devsquadron/cli/constants"
	"github.com/devsquadron/cli/exception"

	"github.com/devsquadron/cli/system"

	"github.com/nanvenomous/exfs"
	"github.com/nanvenomous/snek"
	"github.com/spf13/cobra"
)

var (
	pkgName         = "ds"
	Sys             system.SystemType
	Cfg             configuration.ConfigurationType
	TaskClient      *requests.TaskClient
	DeveloperClient *requests.DeveloperClient
	TeamClient      *requests.TeamClient
	fs              *exfs.FileSystem
	cfgFile         string
	DebugFlag       bool
	version         bool
)

var rootCmd = &cobra.Command{
	Use:   pkgName,
	Short: "devsquadron cli",
	Long:  `devsquadron cli`,
	Run: func(cmd *cobra.Command, args []string) {
		snek.RunShellCompletion(cmd)

		if version {
			fmt.Println(constants.Version)
			os.Exit(0)
		}

		if DebugFlag {
			Cfg.PrintAllPaths()
		}
		cmd.Help()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return Cfg.Write()
	},
}

func Execute() {
	exception.CheckErr(rootCmd.Execute())
}

func init() {
	snek.InitRoot(rootCmd, cfgFile, "devsquadron")

	rootCmd.Flags().BoolVarP(
		&version, "version", "v", false,
		"show the version of this binary",
	)

	rootCmd.Flags().BoolVarP(
		&DebugFlag, "debug", "d", false,
		"print debug output",
	)

	Sys, Cfg, TaskClient, DeveloperClient, TeamClient = system.GlobalSetup()
	cobra.OnInitialize(func() {
		Cfg.Setup(cfgFile)
	})
}
