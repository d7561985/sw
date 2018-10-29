package cmd

import (
	"fmt"
	"github.com/d7561985/sw/generators/newapp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var GenerateCMD = &cobra.Command{
	Use:     "generate",
	Short:   "Swagger generators",
	Aliases: []string{"g"},
	RunE: func(cmd *cobra.Command, args []string) error {
		currpath, _ := os.Getwd()

		logrus.Infof("current path: %s", currpath)

		r, err := cmd.Flags().GetString("framework")
		fmt.Println(r)

		app, err := newapp.New(r)
		if err != nil {
			return err
		}

		return app.Run(currpath)
	},
}

func init() {
	RootCmd.AddCommand(GenerateCMD)
	GenerateCMD.Flags().StringP("framework", "f", newapp.GeneratorList[0], fmt.Sprintf("generate swagger documentation accoirding FrameWork [%s]", strings.Join(newapp.GeneratorList, " ,")))
}
