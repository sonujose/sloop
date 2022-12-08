package template

import "github.com/sirupsen/logrus"

type Template struct {
}

func (t *Template) GeneratePackageTemplates(output string, log *logrus.Logger) error {

	log.Infof("Executing sloop template to %s output", output)

	return nil
}
