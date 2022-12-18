package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	"gopkg.in/yaml.v2"
)

const (
	defaultSloopPkgFile string = "./sloop.yaml"
)

func ParseSloopPackage(log *logrus.Logger, sloopPkgFile string) (*client.SloopPackage, error) {

	if sloopPkgFile == "" {
		sloopPkgFile = defaultSloopPkgFile
	}

	// Read the sloop pkg yaml file and parse the details
	filename, _ := filepath.Abs(sloopPkgFile)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Errorf("Unable to find the sloop package file, please check if the file path is correct.")
		return nil, err
	}

	var sloopPkg client.SloopPackage

	err = yaml.Unmarshal(yamlFile, &sloopPkg)
	if err != nil {
		log.Errorf("Error parsing the sloop package")
		return nil, err
	}

	return &sloopPkg, nil
}
