package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Aliases: []string{"dash"},
	Short:   "Opens the sloop dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		sloopPkg, err := config.ParseSloopPackage(l, packagefile)
		if err != nil {
			return err
		}

		// TODO : Implement logic
		l.Info(sloopPkg)

		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
