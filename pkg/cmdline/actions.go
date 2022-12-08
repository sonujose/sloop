package cmdline

import (
	"github.com/sonujose/sloop/pkg/core/template"
	"github.com/sonujose/sloop/pkg/logger"
	"github.com/urfave/cli/v2"
)

func sloopTemplateAction(ctx *cli.Context) error {

	log := logger.NewLogger()

	tr := template.Template{}
	err := tr.GeneratePackageTemplates(ctx.String("output"), log)

	if err != nil {
		return cli.Exit(err.Error(), 2)
	}

	return nil
}
