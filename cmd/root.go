package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sloop",
	Short: "Sloop is a package manager for lazy kubernetes folks",
	Long:  `A Fast and extensible kubernetes package manager for application deployment and management.`,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

var (
	loglevel   uint32
	configFile string
)

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().Uint32VarP(&loglevel, "verbosity", "v", uint32(logrus.InfoLevel), "verbosity level of the log")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "f", "", "sloop config file path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
