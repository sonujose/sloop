package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	"github.com/sonujose/sloop/apis/v1/controller"
	"github.com/sonujose/sloop/pkg/core/consts"
	"github.com/sonujose/sloop/pkg/core/history"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SyncConfig struct {
	sloopPkg      *client.SloopPackage
	controllerCfg *controller.SloopControllerConfig
	kubeClient    *kubernetes.Clientset
}

func New(cfg *client.SloopPackage, ctrlCfg *controller.SloopControllerConfig, kclient *kubernetes.Clientset) *SyncConfig {
	return &SyncConfig{sloopPkg: cfg, controllerCfg: ctrlCfg, kubeClient: kclient}
}

//SyncPackage - Compiles sloop package to generate sloop-config and apply the same to the sloop namespace
func (s *SyncConfig) SyncPackage(l *logrus.Logger) error {

	var lastSyncRevision int

	hs := history.New(s.kubeClient)

	labelSel := history.GetPackageHistoryLabelSelectors(history.PackageHistoryFilterKey, s.sloopPkg.Metadata.Name)
	syncHistory, err := hs.GetPackageSyncHistorybyLabels(labelSel, s.sloopPkg.Metadata.Namespace)

	if err != nil || len(syncHistory) == 0 {
		lastSyncRevision = 0
	} else {
		lastSyncRevision, err = getLastSyncRevision(syncHistory)

		if err != nil {
			l.Debugf("Unable to find any sync revison for the specified sloop package, creating the first revision. reason=%v", err)
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
		l.Errorf("Error creating sloop config secret.")
		return err
	}

	sloopDeploymentStatus, _ := yaml.Marshal(s.controllerCfg.Status)

	// CONSOLE-INFO : SYNC OPERATION STATUS
	fmt.Println("Succeeded!! Registered sloop configuration for the package")
	fmt.Println(string(sloopDeploymentStatus))

	return nil
}

/*
generateSloopConfigCurrentSecretName
sconfig.sloop.v100.sloop-config-dev.v2
config.sloop.v<version>.<sloop-configname>.v<rev.
*/
func (s *SyncConfig) getSloopConfigSecretName(syncRevision int) string {
	sloopVersionStr := strings.Replace(s.sloopPkg.Spec.Version, ".", "", -1)
	return fmt.Sprintf("config.sloop.v%s.%s.v%d", sloopVersionStr, s.sloopPkg.Metadata.Name, syncRevision)
}

// getSloopConfigSecretObj
// Returns the Kubernetes Secret object for sloop-config
func (s *SyncConfig) getSloopConfigSecretObj(secretName string, secretData []byte, status controller.SloopConfigStatus) *corev1.Secret {

	kcSecretBlob := make(map[string][]byte)
	syncLabels := &consts.SloopConfigSecretLabels{
		Components: fmt.Sprintf("0_%d", len(s.sloopPkg.Spec.Components)),
		Package:    s.sloopPkg.Metadata.Name,
		Revision:   fmt.Sprint(status.SyncRevision),
		Status:     consts.StatusRegistered,
		Updated:    fmt.Sprint(status.DeployedOn.Unix()),
		Version:    s.sloopPkg.Spec.Version,
		Owner:      consts.ToolName,
	}

	kcSecretBlob["config"] = secretData

	SloopControllerSecret := &corev1.Secret{
		Data: kcSecretBlob,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: s.sloopPkg.Metadata.Namespace,
			Labels:    syncLabels.GetConfigLabels(),
		},
		TypeMeta: metav1.TypeMeta{},
		Type:     corev1.SecretType(fmt.Sprintf("sloop.io/%s", s.sloopPkg.Metadata.Name)),
	}

	return SloopControllerSecret

}
