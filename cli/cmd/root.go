/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"github.com/devsquadron/requests"

	"github.com/devsquadron/ds/configuration"

	"github.com/devsquadron/ds/system"

	"github.com/nanvenomous/exfs"
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
)

var rootCmd = &cobra.Command{
	Use:   pkgName,
	Short: "devsquadron cli",
	Long:  `devsquadron cli`,
	Run: func(cmd *cobra.Command, args []string) {
		system.RunShellCompletion(cmd)
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
	system.RootExecution(rootCmd)
}

func init() {
	system.InitRoot(rootCmd, cfgFile)

	rootCmd.Flags().BoolVarP(
		&DebugFlag, "debug", "d", false,
		"print debug output",
	)

	Sys, Cfg, TaskClient, DeveloperClient, TeamClient = system.GlobalSetup()
	cobra.OnInitialize(func() {
		Cfg.Setup(cfgFile)
	})
}
