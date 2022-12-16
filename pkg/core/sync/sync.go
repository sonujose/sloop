package sync

import (
	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/config"
	"github.com/sonujose/sloop/apis/v1/controller"
)

type SyncConfig struct {
	SloopCfg         *config.SloopConfig
	TemplateManifest *string
}

func New(cfg *config.SloopConfig, manifest *string) *SyncConfig {
	return &SyncConfig{SloopCfg: cfg, TemplateManifest: manifest}
}

func (s *SyncConfig) syncPackage(l *logrus.Logger) {

	var controllerConfig = &controller.SloopControllerConfig{
		Name: s.SloopCfg.Metadata.Name,
	}

	l.Info(controllerConfig)
}

func (s *SyncConfig) generateConfigForController() {

}

func (s *SyncConfig) syncSloopControllerConfig() {

}

func (s *SyncConfig) CreateNewConfigForControllerSync() {

}
