package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/config"
	"github.com/sonujose/sloop/apis/v1/controller"
	"github.com/sonujose/sloop/pkg/core/history"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SyncConfig struct {
	SloopCfg      *config.SloopConfig
	ControllerCfg *controller.SloopControllerConfig
	KubeClient    *kubernetes.Clientset
}

func New(cfg *config.SloopConfig, ctrlCfg *controller.SloopControllerConfig, kclient *kubernetes.Clientset) *SyncConfig {
	return &SyncConfig{SloopCfg: cfg, ControllerCfg: ctrlCfg, KubeClient: kclient}
}

func (s *SyncConfig) SyncPackage(l *logrus.Logger) error {

	kcSecretBlob := make(map[string][]byte)

	deployedOn := time.Now()
	var lastSyncRevision int

	lastSyncRevision, err := getLastSyncRevision(s.KubeClient, s.SloopCfg.Metadata.Name, "sloop")
	if err != nil {
		l.Debugf("Unable to find any sync revison for the specified sloop package, creating the first revision. reason=%v", err)
	}

	newSyncRevision := lastSyncRevision + 1

	configStatus := controller.Status{
		DeployedOn:   fmt.Sprint(deployedOn),
		Version:      s.SloopCfg.Spec.Version,
		SyncRevision: newSyncRevision,
		Components:   len(s.SloopCfg.Spec.Components),
		Name:         s.SloopCfg.Metadata.Name,
	}

	s.ControllerCfg.Status = configStatus
	blob, _ := json.Marshal(s.ControllerCfg)
	kcSecretBlob["config"] = blob

	secretName := getSloopConfigSecretName(s.SloopCfg.Spec.Version, s.SloopCfg.Metadata.Name, newSyncRevision)

	sloopConfigSecretLabels := map[string]string{
		"modifiedAt": fmt.Sprint(deployedOn.Unix()),
		"name":       s.SloopCfg.Metadata.Name,
		"owner":      "sloop",
		"revision":   fmt.Sprint(newSyncRevision),
	}

	SloopControllerSecret := &corev1.Secret{
		Data: kcSecretBlob,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "sloop",
			Labels:    sloopConfigSecretLabels,
		},
		TypeMeta: metav1.TypeMeta{},
		Type:     corev1.SecretType(fmt.Sprintf("sloop.io/%s", s.SloopCfg.Metadata.Name)),
	}

	ctx := context.Background()
	_, err = s.KubeClient.CoreV1().Secrets("sloop").Create(ctx, SloopControllerSecret, metav1.CreateOptions{})

	if err != nil {
		l.Errorf("Error creating sloop config secret.")
		return err
	}

	sloopDeploymentStatus, _ := yaml.Marshal(s.ControllerCfg.Status)

	// CONSOLE-INFO : SYNC OPERATION STATUS
	fmt.Println("synced sloop configuration")
	fmt.Println(string(sloopDeploymentStatus))

	return nil
}

func getLastSyncRevision(KubeClient *kubernetes.Clientset, name string, namespace string) (int, error) {
	hs := history.New(KubeClient)
	syncHistory, err := hs.GetPackageSyncHistory(name, namespace)

	var lastrev int = 0

	if err != nil {
		log.Println("Unable to list package history", err)
		return lastrev, fmt.Errorf("No previous sync revisions found")
	}

	var lastSyncExecutionMeta v1.Secret

	if len(syncHistory) >= 0 {
		lastSyncExecutionMeta = syncHistory[0]
	} else {
		return lastrev, fmt.Errorf("No previous sync revisions found")
	}

	lastrev, err = strconv.Atoi(lastSyncExecutionMeta.Labels["revision"])
	if err != nil {
		return lastrev, err
	}

	return lastrev, nil
}

/*
generateSloopConfigCurrentSecretName
sconfig.sloop.v100.sloop-config-dev.v2
config.sloop.v<version>.<sloop-configname>.v<rev.
*/
func getSloopConfigSecretName(version string, name string, syncVersion int) string {
	sloopVersionStr := strings.Replace(version, ".", "", -1)
	return fmt.Sprintf("config.sloop.v%s.%s.v%d", sloopVersionStr, name, syncVersion)
}
