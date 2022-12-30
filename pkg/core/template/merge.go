package template

import (
	"io/ioutil"
	"strings"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	yaml "gopkg.in/yaml.v2"
)

func mergeAllValuesFile(valuesFiles []string, log *logrus.Logger) (map[string]interface{}, error) {

	defValuesFile, err := ioutil.ReadFile(valuesFiles[0])
	if err != nil {
		log.Errorf("Error loading file %s", valuesFiles[0])
		return nil, err
	}

	var defvaluesFileYaml map[string]interface{}
	err = yaml.Unmarshal(defValuesFile, &defvaluesFileYaml)
	if err != nil {
		log.Errorf("Error unmarshall default values file")
		return nil, err
	}

	for _, j := range valuesFiles[1:] {

		valuesFile, err := ioutil.ReadFile(j)

		if err != nil {
			log.Errorf("Error loading values file from sloop components")
			return nil, err
		}

		var valuesFileYaml map[string]interface{}
		err = yaml.Unmarshal(valuesFile, &valuesFileYaml)
		if err != nil {
			log.Errorf("Error unmarshall default values file")
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

func generateValuesForOverrides(keys []string, val string) map[string]interface{} {

	finalElement := map[string]interface{}{
		keys[len(keys)-1]: val,
	}

	for i := len(keys) - 2; i > 0; i-- {
		newElement := map[string]interface{}{
			keys[i]: finalElement,
		}
		finalElement = newElement
	}

	return finalElement
}
