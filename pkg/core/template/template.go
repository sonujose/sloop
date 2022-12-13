package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/configurator"
	"github.com/sonujose/sloop/pkg/logger"
	yaml "gopkg.in/yaml.v2"
)

type Template struct {
	Logger *logrus.Logger
}

const (
	defaultSloopConfigFileName string = "./sloop.yaml"
)

func (t *Template) ParseSloopConfig() *configurator.SloopConfigurator {
	// Read the sloop condif yaml file and parse the details
	filename, _ := filepath.Abs(defaultSloopConfigFileName)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		t.Logger.WithError(err).Errorf("Error unmarshall system")
		t.Logger.Exit(logger.SloopConfiguratorFileLoadError)
	}

	var config configurator.SloopConfigurator

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		t.Logger.WithError(err).Errorf("Error parsing the sloop config")
		t.Logger.Exit(logger.SloopConfiguratorFileParseError)
	}

	fmt.Printf("Value: %#v\n", config.Kind)

	return &config
}

func (t *Template) GeneratePackageTemplates() error {

	configuratorObj := t.ParseSloopConfig()

	for _, j := range configuratorObj.Spec.Components {

		if !j.Enabled {
			continue
		}

		templatePath := j.Path

		templateFiles, err := walkMatch(templatePath, "*")

		if err != nil {
			t.Logger.WithError(err).Errorf("Error getting path....")
			t.Logger.Exit(1)
		}

		t.Logger.Info(templateFiles)

		valuesFile, err := ioutil.ReadFile(j.ValuesFiles[0])

		if err != nil {
			t.Logger.WithError(err).Errorf("Error loading values file fr component")
			t.Logger.Exit(1)
		}

		var componentManifest bytes.Buffer

		for _, templateFile := range templateFiles {
			data, err := ioutil.ReadFile(templateFile)
			if err != nil {
				t.Logger.WithError(err).Errorf("Error loading file")
				t.Logger.Exit(1)
			}

			blobString := string(data)

			var d2 map[string]interface{}
			err = yaml.Unmarshal(valuesFile, &d2)

			overrideValues := j.Set

			var overridesYamls map[string]interface{}

			for _, n := range overrideValues {
				mapKey := n.Name

				//var dbs map[string]interface{}

				setKeys := strings.Split(mapKey, ".")

				newmap := generateValuesForOverrides(setKeys, n.Value)
				//t.Logger.Info(newmap)

				overridesYamls = mergeMaps(overridesYamls, newmap)

			}

			t.Logger.Infof("Override yamls - %v", overridesYamls)
			t.Logger.Infof("Old base yaml - %v", d2)

			newUpdatedYaml := mergeMaps(d2, overridesYamls)

			t.Logger.Infof("New updated yaml - %v", newUpdatedYaml)

			if err != nil {
				t.Logger.WithError(err).Errorf("Error unmarshall values file")
				t.Logger.Exit(1)
			}

			componentManifest.Write([]byte("---\n"))
			componentManifest.Write([]byte(fmt.Sprintf("# Component: %s, Manifest: %s\n", j.Name, templateFile)))
			//t.Logger.Info(d2)
			err = tpl(blobString, newUpdatedYaml, &componentManifest)

			if err != nil {
				t.Logger.WithError(err).Errorf("Failed to generate templates for the manifests")
				t.Logger.Exit(1)
			}

		}

		fmt.Print(componentManifest.String())

	}

	return nil
}
