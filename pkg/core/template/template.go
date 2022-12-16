package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"io"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/config"
	"github.com/sonujose/sloop/apis/v1/controller"
	yaml "gopkg.in/yaml.v2"
)

type SloopTemplate struct {
	SloopCfg *config.SloopConfig
}

func New(cfg *config.SloopConfig) *SloopTemplate {
	return &SloopTemplate{SloopCfg: cfg}
}

func (t *SloopTemplate) GeneratePackageTemplates(log *logrus.Logger) (*controller.SloopControllerConfig, error) {

	configuratorObj := t.SloopCfg

	var componentManifest bytes.Buffer

	sloopSyncConfiguration := &controller.SloopControllerConfig{
		Name: t.SloopCfg.Metadata.Name,
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

		valuesFile, err := ioutil.ReadFile(j.ValuesFiles[0])

		if err != nil {
			log.WithError(err).Errorf("Error loading values file from sloop components")
			return nil, err
		}

		componentConfig := &controller.Component{
			Name:      j.Name,
			Namespace: j.Namespace,
			Path:      j.Path,
		}

		for _, templateFile := range templateFiles {
			data, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.WithError(err).Errorf("Error loading template file - %s", templateFile)
				return nil, err
			}

			blobString := string(data)

			var valuesFileYaml map[string]interface{}
			err = yaml.Unmarshal(valuesFile, &valuesFileYaml)
			if err != nil {
				log.WithError(err).Errorf("Error unmarshall default values file")
				return nil, err
			}

			overrideValues := j.Set

			var overridesYaml map[string]interface{}

			for _, n := range overrideValues {
				mapKey := n.Name
				setKeys := strings.Split(mapKey, ".")
				newmap := generateValuesForOverrides(setKeys, n.Value)
				overridesYaml = mergeMaps(overridesYaml, newmap)
			}

			updatedMergedYaml := mergeMaps(valuesFileYaml, overridesYaml)

			// DEBUG
			//log.Debugf("Default values yaml - %v", valuesFileYaml)
			//valuesFile, _ := yaml.Marshal(valuesFileYaml)
			//log.Infof("Default values yaml\n%s", string(valuesFile))

			//log.Debugf("Override yaml - %v", overridesYaml)
			//overrideYamlFile, _ := yaml.Marshal(overridesYaml)
			//log.Infof("Override Yaml\n%s", string(overrideYamlFile))

			//log.Debugf("New updated yaml - %v", updatedMergedYaml)
			//mergedYamlFile, _ := yaml.Marshal(updatedMergedYaml)
			//log.Infof("Merged Yaml\n%s", string(mergedYamlFile))

			componentManifest.Write([]byte("---\n"))
			componentManifest.Write([]byte(fmt.Sprintf("# Component: %s, Manifest: %s\n", j.Name, templateFile)))

			err = templateExecutor(blobString, updatedMergedYaml, &componentManifest)

			if err != nil {
				log.WithError(err).Errorf("Failed to generate templates for the manifests")
				return nil, err
			}

			templateMeta := &controller.TemplateFile{
				FileName: templateFile,
			}

			componentConfig.TemplateFiles = append(componentConfig.TemplateFiles, *templateMeta)
		}

		sloopSyncConfiguration.Config.Components = append(sloopSyncConfiguration.Config.Components, *componentConfig)
	}

	sloopSyncConfiguration.Config.ConsolidatedManifest = componentManifest.String()

	return sloopSyncConfiguration, nil
}

func templateExecutor(t string, vals map[string]interface{}, out io.Writer) error {
	tt, err := template.New("_").Parse(t)
	if err != nil {
		return err
	}
	return tt.Execute(out, vals)
}
