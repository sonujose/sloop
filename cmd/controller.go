package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:     "controller",
	Aliases: []string{"viz"},
	Short:   "Setup the sloop controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		sloopConfig, err := config.ParseSloopConfig(l, configFile)
		if err != nil {
			return err
		}

		// TODO : Implement logic to analyze and generate report for sloop
		l.Info(sloopConfig)

		return nil
	},
	SilenceUsage: true,
}

var controllerInstallCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "install the sloop controller in the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		sloopConfig, err := config.ParseSloopConfig(l, configFile)
		if err != nil {
			return err
		}

		// TODO : Implement logic to analyze and generate report for sloop
		l.Info(sloopConfig)

		return nil
	},
	SilenceUsage: true,
}

func init() {
	controllerCmd.AddCommand(controllerInstallCmd)
	rootCmd.AddCommand(controllerCmd)
}
