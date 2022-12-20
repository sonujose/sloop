package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	"github.com/sonujose/sloop/apis/v1/controller"
	"github.com/sonujose/sloop/pkg/core/history"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SyncConfig struct {
	sloopPkg      *client.SloopPackage
	controllerCfg *controller.SloopControllerConfig
	kubeClient    *kubernetes.Clientset
	log           *logrus.Logger
}

func New(cfg *client.SloopPackage, ctrlCfg *controller.SloopControllerConfig, kclient *kubernetes.Clientset, l *logrus.Logger) *SyncConfig {
	return &SyncConfig{sloopPkg: cfg, controllerCfg: ctrlCfg, kubeClient: kclient, log: l}
}

//SyncPackage - Compiles sloop package to generate sloop-config and apply the same to the sloop namespace
//
// Executes on the command `sloop sync`
//
func (s *SyncConfig) SyncPackage() error {

	var lastSyncRevision int

	hs := history.New(s.kubeClient, s.log)

	labelSel := hs.GetPackageHistoryLabelSelectors(history.PackageHistoryFilterKey, s.sloopPkg.Metadata.Name, "")
	syncHistory, err := hs.GetPackageSyncHistorybyLabels(labelSel, s.sloopPkg.Metadata.Namespace)

	if err != nil || len(syncHistory) == 0 {
		lastSyncRevision = 0
	} else {
		lastSyncRevision, err = getLastSyncRevision(syncHistory)

		if err != nil {
			s.log.Debugf("Unable to find any sync revison for the specified sloop package, creating the first revision. reason=%v", err)
			lastSyncRevision = 0
		}
	}

	syncRevision := lastSyncRevision + 1

	deployedOn := time.Now()

	configStatus := controller.SloopConfigStatus{
		DeployedOn:   deployedOn,
		Version:      s.sloopPkg.Spec.Version,
		SyncRevision: syncRevision,
		Components:   fmt.Sprintf("0/%d", len(s.sloopPkg.Spec.Components)),
		Name:         s.sloopPkg.Metadata.Name,
	}

	s.controllerCfg.Status = configStatus

	secretDataBlob, _ := json.Marshal(s.controllerCfg)
	secretName := s.getSloopConfigSecretName(syncRevision)

	SloopControllerSecret := s.getSloopConfigSecretObj(secretName, secretDataBlob, configStatus)

	ctx := context.Background()
	_, err = s.kubeClient.CoreV1().Secrets(s.sloopPkg.Metadata.Namespace).Create(ctx, SloopControllerSecret, metav1.CreateOptions{})

	if err != nil {
		s.log.Errorf("Error creating sloop config secret.")
		return err
	}

	sloopDeploymentStatus, _ := yaml.Marshal(s.controllerCfg.Status)

	// Aborting pending revisions which are still pending.
	s.abortOldRegisteredRevisions(fmt.Sprint(configStatus.SyncRevision))

	// CONSOLE-INFO : SYNC OPERATION STATUS
	fmt.Println(string(sloopDeploymentStatus))
	fmt.Println("Succeeded!! Registered sloop configuration for the package")

	return nil
}
