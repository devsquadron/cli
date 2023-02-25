/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"github.com/nanvenomous/snek"
	"github.com/spf13/cobra"
)

const (
	PKG_NAME = "pm"
)

var (
	cfgFile string
	prod    bool
)

var rootCmd = &cobra.Command{
	Use:   PKG_NAME,
	Short: "some development help for this project",
	Long:  `some development help for this project`,
	Run: func(cmd *cobra.Command, args []string) {
		snek.RunShellCompletion(cmd)
		cmd.Help()
	},
}

func Execute() {
	var err error
	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	snek.InitRoot(rootCmd, cfgFile, "devsquadron")
	rootCmd.PersistentFlags().BoolVarP(&prod, "production", "p", false, "run with production mongo instance.")
}
