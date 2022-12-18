package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:     "controller",
	Aliases: []string{"c"},
	Short:   "sloop controller operations",
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

var controllerInstallCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i"},
	Short:   "install the sloop controller in the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		sloopPkg, err := config.ParseSloopPackage(l, packagefile)
		if err != nil {
			return err
		}

		// TODO : Implement logic to install sloop controller
		l.Info(sloopPkg)

		return nil
	},
	SilenceUsage: true,
}

func init() {
	controllerCmd.AddCommand(controllerInstallCmd)
	rootCmd.AddCommand(controllerCmd)
}
