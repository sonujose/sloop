package cmd

import (
	"os"

	"github.com/sonujose/sloop/pkg/cmdline"
)

func Execute() {

	app := cmdline.Initialize()

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
