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
	revision    string
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Creates the desired state from the sloop package file and apply the sloop config to the cluster",
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

		sloopCtrlInput, err := json.Marshal(sloopCtrlConfig)
		if err != nil {
			return err
		}

		l.Tracef("Sloop Controller config - %v", string(sloopCtrlInput))

		kclient, err := kube.NewClient()

		if err != nil {
			l.Errorf("Error connecting to cluster via clientcmd. Aborting...")
			return err
		}

		ss := sloopSync.New(sloopPkg, sloopCtrlConfig, kclient, l)

		err = ss.SyncPackage()

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

		sloopPkg, err := config.ParseSloopPackage(l, packagefile)
		if err != nil {
			return err
		}

		kclient, err := kube.NewClient()

		if err != nil {
			l.Errorf("Error connecting to cluster via clientcmd. Aborting...")
			return err
		}

		hs := history.New(kclient, l)

		syncHistory, err := hs.ListSyncHistory(sloopPkg.Metadata.Name, sloopPkg.Metadata.Namespace, allpackages)
		if err != nil {
			return err
		}

		if len(syncHistory) > 0 {
			console.PrintSyncHistory(syncHistory)
		} else {
			fmt.Println("No sync history found for the package -", sloopPkg.Metadata.Name)
		}

		return nil
	},
	SilenceUsage: true,
}

var syncCleanCmd = &cobra.Command{
	Use:     "clean",
	Aliases: []string{"c"},
	Short:   "Cleans the specified revision metadata of the package from the sloop db (Use 'delete' command to remove the components installation from the cluster.)",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		return nil
	},
	SilenceUsage: true,
}

var syncDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "Use this command to remove the components installation from the cluster.)",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.NewLogger()

		l.SetLevel(logrus.Level(loglevel))

		return nil
	},
	SilenceUsage: true,
}

func init() {
	syncCmd.PersistentFlags().BoolVarP(&dryrun, "dryrun", "d", false, "Use dry run to test the sloop configurations.")
	syncCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "specify the namespace for syncing sloop config")
	syncHistoryCmd.PersistentFlags().BoolVarP(&allpackages, "allpackages", "a", false, "get history from all the packages")
	syncCleanCmd.PersistentFlags().StringVarP(&revision, "revision", "r", "", "specify the sync revision that needs to be cleaned")
	syncCmd.AddCommand(syncHistoryCmd)
	syncCmd.AddCommand(syncCleanCmd)
	syncCmd.AddCommand(syncDeleteCmd)
	rootCmd.AddCommand(syncCmd)
}
