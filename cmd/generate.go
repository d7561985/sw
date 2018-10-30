package cmd

import (
	"github.com/d7561985/sw/generators/swaggergen"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var GenerateCMD = &cobra.Command{
	Use:     "generate",
	Short:   "Swagger generators",
	Aliases: []string{"g"},
	RunE: func(cmd *cobra.Command, args []string) error {
		currpath, _ := os.Getwd()

		logrus.Infof("current path: %s", currpath)

		swaggergen.GenerateDocs(currpath)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(GenerateCMD)
}
