package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/config"
	yaml "gopkg.in/yaml.v2"
)

type SloopTemplate struct {
	SloopCfg *config.SloopConfig
}

func New(cfg *config.SloopConfig) *SloopTemplate {
	return &SloopTemplate{SloopCfg: cfg}
}

func (t *SloopTemplate) GeneratePackageTemplates(log *logrus.Logger) error {

	configuratorObj := t.SloopCfg

	for _, j := range configuratorObj.Spec.Components {

		if !j.Enabled {
			continue
		}

		templatePath := j.Path

		templateFiles, err := walkMatch(templatePath, "*")

		if err != nil {
			log.WithError(err).Errorf("Error getting path....")
			log.Exit(1)
		}

		log.Debug(templateFiles)

		valuesFile, err := ioutil.ReadFile(j.ValuesFiles[0])

		if err != nil {
			log.WithError(err).Errorf("Error loading values file from sloop components")
			return err
		}

		var componentManifest bytes.Buffer

		for _, templateFile := range templateFiles {
			data, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.WithError(err).Errorf("Error loading template file")
				return err
			}

			blobString := string(data)

			var valuesFileYaml map[string]interface{}
			err = yaml.Unmarshal(valuesFile, &valuesFileYaml)
			if err != nil {
				log.WithError(err).Errorf("Error unmarshall values file")
				return err
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
			valuesFile, _ := yaml.Marshal(valuesFileYaml)
			log.Infof("\n%s", string(valuesFile))

			//log.Debugf("Override yaml - %v", overridesYaml)
			overrideYamlFile, _ := yaml.Marshal(overridesYaml)
			log.Infof("\n%s", string(overrideYamlFile))

			//log.Debugf("New updated yaml - %v", updatedMergedYaml)
			mergedYamlFile, _ := yaml.Marshal(updatedMergedYaml)
			log.Infof("\n%s", string(mergedYamlFile))

			componentManifest.Write([]byte("---\n"))
			componentManifest.Write([]byte(fmt.Sprintf("# Component: %s, Manifest: %s\n", j.Name, templateFile)))

			err = tpl(blobString, updatedMergedYaml, &componentManifest)

			if err != nil {
				log.WithError(err).Errorf("Failed to generate templates for the manifests")
				return err
			}

		}

		log.Infof("\n%s", componentManifest.String())

	}

	return nil
}
