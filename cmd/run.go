/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/devsquadron/project-manager/constants"
	"github.com/devsquadron/project-manager/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var (
	quiet bool
	mock  bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "start the service to listen on 5000",
	Long:  `start the service to listen on 5000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err  error
			rtr  *gin.Engine
			clnt *mongo.Client
			ctx  = context.Background()
		)

		rtr, clnt, _, err = services.GetService(ctx, prod, quiet)
		if err != nil {
			return err
		}
		fmt.Println("Got Service")

		err = rtr.Run(constants.PORT)
		if err != nil {
			return err
		}

		return clnt.Disconnect(ctx)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "run in quiet mode. log to file")
	runCmd.PersistentFlags().BoolVarP(&mock, "mock", "m", false, "run with mock data.")
}
