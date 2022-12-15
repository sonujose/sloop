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
			log.WithError(err).Errorf("Error loading values file fr component")
			return err
		}

		var componentManifest bytes.Buffer

		for _, templateFile := range templateFiles {
			data, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.WithError(err).Errorf("Error loading file")

				return err
			}

			blobString := string(data)

			var d2 map[string]interface{}
			err = yaml.Unmarshal(valuesFile, &d2)

			overrideValues := j.Set

			var overridesYamls map[string]interface{}

			for _, n := range overrideValues {
				mapKey := n.Name

				setKeys := strings.Split(mapKey, ".")

				newmap := generateValuesForOverrides(setKeys, n.Value)

				overridesYamls = mergeMaps(overridesYamls, newmap)

			}

			log.Debugf("Override yamls - %v", overridesYamls)
			log.Debugf("Old base yaml - %v", d2)

			newUpdatedYaml := mergeMaps(d2, overridesYamls)

			log.Debugf("New updated yaml - %v", newUpdatedYaml)

			if err != nil {
				log.WithError(err).Errorf("Error unmarshall values file")
				return err
			}

			componentManifest.Write([]byte("---\n"))
			componentManifest.Write([]byte(fmt.Sprintf("# Component: %s, Manifest: %s\n", j.Name, templateFile)))

			err = tpl(blobString, newUpdatedYaml, &componentManifest)

			if err != nil {
				log.WithError(err).Errorf("Failed to generate templates for the manifests")
				return err
			}

		}

		log.Info(componentManifest.String())

	}

	return nil
}
