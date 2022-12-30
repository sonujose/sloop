package template

import (
	"bytes"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	"github.com/sonujose/sloop/apis/v1/controller"
)

type SloopTemplate struct {
	sloopPkg *client.SloopPackage
}

func New(pkg *client.SloopPackage) *SloopTemplate {
	return &SloopTemplate{sloopPkg: pkg}
}

func (t *SloopTemplate) GeneratePackageTemplates(log *logrus.Logger) (*controller.SloopControllerConfig, error) {

	configuratorObj := t.sloopPkg
	var componentManifest bytes.Buffer

	sloopSyncConfiguration := &controller.SloopControllerConfig{
		Name: t.sloopPkg.Metadata.Name,
	}

	for _, comp := range configuratorObj.Spec.Components {

		if !comp.Enabled {
			continue
		}

		templatePath := comp.Path
		templateFiles, err := walkMatch(templatePath, "*")

		if err != nil {
			log.WithError(err).Errorf("Error getting template files from the specified path %s....", templatePath)
			return nil, err
		}

		log.Debug(templateFiles)

		mergedValuesFileyaml, err := coalesceAllValuesFiles(comp.ValuesFiles, log)
		if err != nil {
			return nil, err
		}

		updatedMergedYaml := coalesceValuesFilesWithOverrides(comp.Set, mergedValuesFileyaml)

		componentConfig := &controller.Component{
			Name:      comp.Name,
			Namespace: comp.Namespace,
			Path:      comp.Path,
		}

		for _, templateFile := range templateFiles {

			var templateManifest bytes.Buffer

			templateFileBlob, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.Errorf("Error loading template file - %s", templateFile)
				return nil, err
			}

			err = templateExecutor(templateFileBlob, updatedMergedYaml, &templateManifest)

			if err != nil {
				log.Errorf("Failed to generate templates for the manifests")
				return nil, err
			}

			templateMeta := &controller.TemplateFile{
				FileName:     templateFile,
				ManifestYaml: templateManifest.String(),
			}

			componentConfig.TemplateFiles = append(componentConfig.TemplateFiles, *templateMeta)
			templateManifestWithMeta := appendTemplateManifestWithMeta(templateManifest, comp, templateFile)
			componentManifest.Write([]byte(templateManifestWithMeta.String()))

		}

		sloopSyncConfiguration.Config.Components = append(sloopSyncConfiguration.Config.Components, *componentConfig)
	}

	sloopSyncConfiguration.Config.ConsolidatedManifest = componentManifest.String()

	return sloopSyncConfiguration, nil
}
