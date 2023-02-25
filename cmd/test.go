/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/devsquadron/project-manager/database"
	"github.com/devsquadron/project-manager/security"
	"github.com/devsquadron/project-manager/services"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test db ops",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err   error
			clnt  *mongo.Client
			ctx   = context.Background()
			tmDtb *database.TeamDatabase
		)

		clnt, _, _, _, tmDtb, err = services.GetDatabases(ctx, prod, quiet, security.NewSecurity())
		if err != nil {
			return err
		}

		info, err := tmDtb.Info("devsquadron")
		if err != nil {
			return err
		}
		fmt.Println(info)

		return clnt.Disconnect(ctx)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
