package template

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"io"
	"text/template"

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

	for _, j := range configuratorObj.Spec.Components {

		if !j.Enabled {
			continue
		}

		templatePath := j.Path
		templateFiles, err := walkMatch(templatePath, "*")

		if err != nil {
			log.WithError(err).Errorf("Error getting template files from the specified path %s....", templatePath)
			return nil, err
		}

		log.Debug(templateFiles)

		mergedValuesFileyaml, err := mergeAllValuesFile(j.ValuesFiles, log)
		updatedMergedYaml := mergeOverrideWithValuesFiles(j.Set, mergedValuesFileyaml)

		componentConfig := &controller.Component{
			Name:      j.Name,
			Namespace: j.Namespace,
			Path:      j.Path,
		}

		for _, templateFile := range templateFiles {

			var templateManifest bytes.Buffer

			templateFileBlob, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.WithError(err).Errorf("Error loading template file - %s", templateFile)
				return nil, err
			}

			err = templateExecutor(templateFileBlob, updatedMergedYaml, &templateManifest)

			if err != nil {
				log.WithError(err).Errorf("Failed to generate templates for the manifests")
				return nil, err
			}

			templateMeta := &controller.TemplateFile{
				FileName:     templateFile,
				ManifestYaml: templateManifest.String(),
			}

			componentConfig.TemplateFiles = append(componentConfig.TemplateFiles, *templateMeta)

			templateManifestWithMeta := appendTemplateManifestWithMeta(templateManifest, j, templateFile)
			componentManifest.Write([]byte(templateManifestWithMeta.String()))

		}

		sloopSyncConfiguration.Config.Components = append(sloopSyncConfiguration.Config.Components, *componentConfig)
	}

	sloopSyncConfiguration.Config.ConsolidatedManifest = componentManifest.String()

	return sloopSyncConfiguration, nil
}

func appendTemplateManifestWithMeta(templateManifest bytes.Buffer, comp client.Component, templateFile string) bytes.Buffer {

	var packageTemplateManifest bytes.Buffer

	// Write the whole template manifest for all components
	packageTemplateManifest.Write([]byte("---\n"))
	packageTemplateManifest.Write([]byte(fmt.Sprintf("# Component: %s, Manifest: %s, Namespace: %s\n", comp.Name, templateFile, comp.Namespace)))
	packageTemplateManifest.Write(templateManifest.Bytes())

	return packageTemplateManifest
}

func templateExecutor(t []byte, vals map[string]interface{}, out io.Writer) error {
	tt, err := template.New("_").Parse(string(t))
	if err != nil {
		return err
	}
	return tt.Execute(out, vals)
}
