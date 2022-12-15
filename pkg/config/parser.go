package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/config"
	"gopkg.in/yaml.v2"
)

const (
	defaultSloopConfigFileName string = "./sloop.yaml"
)

func ParseSloopConfig(log *logrus.Logger, configFile string) (*config.SloopConfig, error) {

	if configFile == "" {
		configFile = defaultSloopConfigFileName
	}

	// Read the sloop condif yaml file and parse the details
	filename, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.WithError(err).Errorf("Unable to find the sloop config file, please check if the file path is correct.")
		return nil, err
	}

	var config config.SloopConfig

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.WithError(err).Errorf("Error parsing the sloop config")
		return nil, err
	}

	return &config, nil
}
