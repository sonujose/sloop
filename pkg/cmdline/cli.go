package cmdline

import (
	"github.com/urfave/cli/v2"
)

func Initialize() *cli.App {
	app := &cli.App{
		Name:                 "sloop",
		Usage:                "A simple Kubernetes Package manager for all lazy folks",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "filename",
				Aliases: []string{"f"},
				Usage:   "specify the sloop config file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "template",
				Usage:   "locally Render the sloop package configurations",
				Aliases: []string{"t"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "specify the standard output for templating",
						EnvVars: []string{"OUTPUT"},
					},
				},
				Action: sloopTemplateAction,
			},
			{
				Name:    "analyse",
				Usage:   "analyse the sloop config file and generate report",
				Aliases: []string{"anly"},
				Action:  sloopTemplateAction,
			},
			{
				Name:    "visualize",
				Usage:   "open the sloop visualizer dashboard",
				Aliases: []string{"viz"},
				Action:  sloopTemplateAction,
			},
			{
				Name:    "controller",
				Usage:   "manage sloop controller",
				Aliases: []string{"ctrl"},
				Subcommands: []*cli.Command{
					{
						Name:  "install",
						Usage: "install the sloop controller in the cluster, if not present",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "namespace",
								Aliases:  []string{"n"},
								Usage:    "Specify the namespace for controller installation",
								Required: true,
							},
						},
					},
					{
						Name:  "status",
						Usage: "check the status of the sloop controller in the cluster",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "namespace",
								Aliases:  []string{"n"},
								Usage:    "Specify the namespace for controller installation",
								Required: true,
							},
						},
					},
				},
			},
			{
				Name:    "apply",
				Usage:   "apply the slooper configurations to the target cluster",
				Aliases: []string{"a"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dry-run",
						Value: false,
						Usage: "Do a dry run for the installation, checks all configurations",
					},
					&cli.StringFlag{
						Name:    "namespace",
						Aliases: []string{"n"},
						Usage:   "Specify the namespace for the target installation",
					},
				},
			},
		},
	}

	return app
}
