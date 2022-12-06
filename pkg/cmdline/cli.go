package cmdline

import (
	"fmt"

	"github.com/sonujose/sloop/pkg/logger"
	"github.com/sonujose/sloop/pkg/template"
	"github.com/urfave/cli/v2"
)

func Initialize() *cli.App {
	app := &cli.App{
		Name:  "sloop",
		Usage: "Sloop package Manager",
		Action: func(*cli.Context) error {
			fmt.Println("Welcome!!! Sloop package manager is ready...")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "template",
				Usage:   "Render the sloop packages as templates to the output source",
				Aliases: []string{"t", "tpl"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o", "out"},
						Usage:   "specify the standard output for templating",
						EnvVars: []string{"OUTPUT"},
					},
				},
				Action: sloopTemplateAction,
			},
		},
	}

	return app
}

func sloopTemplateAction(ctx *cli.Context) error {

	log := logger.NewLogger()

	err := template.GeneratePackageTemplates(ctx.String("output"), log)

	if err != nil {
		return cli.Exit(err.Error(), 2)
	}

	return nil
}
