package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Aliases: []string{"viz"},
	Short:   "Open the sloop visualization dashboard",
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
	rootCmd.AddCommand(dashboardCmd)
}
