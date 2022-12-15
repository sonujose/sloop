package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	dryrun    bool
	namespace string
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Apply the sloop config to the cluster, sync the desired state from the sloop config file",
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
	syncCmd.PersistentFlags().BoolVarP(&dryrun, "dryrun", "d", false, "Use dry run to test the sloop configurations.")
	syncCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "specify the global namespace for syncing sloop config")

	rootCmd.AddCommand(syncCmd)
}
