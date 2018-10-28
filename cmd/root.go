package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd is the hook for all of the other commands in the buffalo binary.
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "swagger",
	Short:         "Helps you build your Buffalo applications that much easier!",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Errorf("Error: %s\n\n", err)
	}
}

func init() {
	RootCmd.AddCommand()
}
