package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sloop",
	Short: "Sloop is a package manager for all lazy kubernetes folks out there",
	Long:  `Sloop is a kubernetes package manager designed for all lazy devops folks. It helps to simplify deploying, packaging and managing kubernetes artifacts.`,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

var (
	loglevel    uint32
	packagefile string
	component   string
)

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().Uint32VarP(&loglevel, "verbosity", "v", uint32(logrus.InfoLevel), "verbosity level of the log")
	rootCmd.PersistentFlags().StringVarP(&packagefile, "packagefile", "f", "", "sloop package file path")
	rootCmd.PersistentFlags().StringVarP(&component, "component", "c", "", "specify component in the package")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}
