package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:     "analyze",
	Aliases: []string{"a"},
	Short:   "Analyze the sloop package file and generate parse report",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		sloopPkg, err := config.ParseSloopPackage(l, packagefile)
		if err != nil {
			return err
		}

		// TODO : Implement logic to analyze and generate report for sloop
		l.Info(sloopPkg)

		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
