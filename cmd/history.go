package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/config"
	"github.com/sonujose/sloop/pkg/console"
	"github.com/sonujose/sloop/pkg/core/history"
	"github.com/sonujose/sloop/pkg/kube"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
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

func init() {
	historyCmd.PersistentFlags().BoolVarP(&allpackages, "allpackages", "a", false, "get history from all the packages")
	rootCmd.AddCommand(historyCmd)
}
