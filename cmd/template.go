package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/controller"
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

		sloopPkg, err := config.ParseSloopPackage(l, packagefile)
		if err != nil {
			return err
		}
		tr := template.New(sloopPkg)

		var sloopCtrlConfig *controller.SloopControllerConfig

		if sloopCtrlConfig, err = tr.GeneratePackageTemplates(l); err != nil {
			return err
		}

		// CONSOLE_INFO : showing the templated manifest
		fmt.Println(sloopCtrlConfig.Config.ConsolidatedManifest)

		sloopcontrollerContent, err := json.Marshal(sloopCtrlConfig)
		if err != nil {
			return err
		}

		l.Tracef("Sloop Controller config - %v", string(sloopcontrollerContent))

		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
