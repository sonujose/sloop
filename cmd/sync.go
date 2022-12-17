package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/controller"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/console"
	"github.com/sonujose/sloop/pkg/core/history"
	sloopSync "github.com/sonujose/sloop/pkg/core/sync"
	"github.com/sonujose/sloop/pkg/core/template"
	"github.com/sonujose/sloop/pkg/kube"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	dryrun      bool
	namespace   string
	allpackages bool
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
		tr := template.New(sloopConfig)

		var sloopCtrlConfig *controller.SloopControllerConfig

		if sloopCtrlConfig, err = tr.GeneratePackageTemplates(l); err != nil {
			return err
		}

		sloopCtrlInput, err := json.Marshal(sloopCtrlConfig)
		if err != nil {
			return err
		}

		l.Debugf("\nSloop Controller config - %v", string(sloopCtrlInput))

		kclient, err := kube.NewClient()

		if err != nil {
			l.Errorf("Error connecting to cluster via clientcmd. Aborting...")
			return err
		}

		ss := sloopSync.New(sloopConfig, sloopCtrlConfig, kclient)

		err = ss.SyncPackage(l)

		if err != nil {
			return err
		}

		return nil
	},
	SilenceUsage: true,
}

var syncHistoryCmd = &cobra.Command{
	Use:     "history",
	Aliases: []string{"h"},
	Short:   "List the sync history if the specified sloop config.",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		sloopConfig, err := config.ParseSloopConfig(l, configFile)
		if err != nil {
			return err
		}

		kclient, err := kube.NewClient()

		if err != nil {
			l.Errorf("Error connecting to cluster via clientcmd. Aborting...")
			return err
		}

		hs := history.New(kclient)

		syncHistory, err := hs.ListSyncHistory(sloopConfig.Metadata.Name, sloopConfig.Metadata.Namespace, allpackages)
		if err != nil {
			return err
		}

		if len(syncHistory) > 0 {
			console.PrintSyncHistory(syncHistory)
		} else {
			fmt.Println("No sync history found for the package -", sloopConfig.Metadata.Name)
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	syncCmd.PersistentFlags().BoolVarP(&dryrun, "dryrun", "d", false, "Use dry run to test the sloop configurations.")
	syncCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "specify the global namespace for syncing sloop config")
	syncHistoryCmd.PersistentFlags().BoolVarP(&allpackages, "allpackages", "a", false, "get history from all the packages")
	syncCmd.AddCommand(syncHistoryCmd)
	rootCmd.AddCommand(syncCmd)
}
