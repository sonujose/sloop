package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"io"
	"text/template"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	"github.com/sonujose/sloop/apis/v1/controller"
	yaml "gopkg.in/yaml.v2"
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

			data, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.WithError(err).Errorf("Error loading template file - %s", templateFile)
				return nil, err
			}

			blobString := string(data)

			templateManifest.Write([]byte("---\n"))
			templateManifest.Write([]byte(fmt.Sprintf("# Component: %s, Manifest: %s, Namespace: %s\n", j.Name, templateFile, j.Namespace)))

			err = templateExecutor(blobString, updatedMergedYaml, &templateManifest)

			if err != nil {
				log.WithError(err).Errorf("Failed to generate templates for the manifests")
				return nil, err
			}

			templateMeta := &controller.TemplateFile{
				FileName:     templateFile,
				ManifestYaml: templateManifest.String(),
			}

			componentConfig.TemplateFiles = append(componentConfig.TemplateFiles, *templateMeta)
			componentManifest.Write([]byte(templateManifest.String()))
		}

		sloopSyncConfiguration.Config.Components = append(sloopSyncConfiguration.Config.Components, *componentConfig)
	}

	sloopSyncConfiguration.Config.ConsolidatedManifest = componentManifest.String()

	return sloopSyncConfiguration, nil
}

func mergeAllValuesFile(valuesFiles []string, log *logrus.Logger) (map[string]interface{}, error) {

	defValuesFile, err := ioutil.ReadFile(valuesFiles[0])

	var defvaluesFileYaml map[string]interface{}
	err = yaml.Unmarshal(defValuesFile, &defvaluesFileYaml)
	if err != nil {
		log.WithError(err).Errorf("Error unmarshall default values file")
		return nil, err
	}

	for _, j := range valuesFiles[1:] {

		valuesFile, err := ioutil.ReadFile(j)

		if err != nil {
			log.WithError(err).Errorf("Error loading values file from sloop components")
			return nil, err
		}

		var valuesFileYaml map[string]interface{}
		err = yaml.Unmarshal(valuesFile, &valuesFileYaml)
		if err != nil {
			log.WithError(err).Errorf("Error unmarshall default values file")
			return nil, err
		}

		err = mergo.Merge(&defvaluesFileYaml, valuesFileYaml, mergo.WithOverride)

		if err != nil {
			return nil, err
		}

	}

	return defvaluesFileYaml, nil
}

func mergeOverrideWithValuesFiles(overrideValues []client.Set, mergedValuesFile map[string]interface{}) map[string]interface{} {

	for _, n := range overrideValues {
		mapKey := n.Name
		setKeys := strings.Split(mapKey, ".")
		newval := generateValuesForOverrides(setKeys, n.Value)
		err := mergo.Merge(&mergedValuesFile, newval, mergo.WithOverride)
		if err != nil {
			continue
		}
	}

	return mergedValuesFile
}

func templateExecutor(t string, vals map[string]interface{}, out io.Writer) error {
	tt, err := template.New("_").Parse(t)
	if err != nil {
		return err
	}
	return tt.Execute(out, vals)
}
