package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/core/template"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:     "template",
	Aliases: []string{"t"},
	Short:   "Locally render the sloop package",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		sloopConfig, err := config.ParseSloopConfig(l, configFile)
		if err != nil {
			return err
		}
		tr := template.New(sloopConfig)

		if err := tr.GeneratePackageTemplates(l); err != nil {
			return err
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
