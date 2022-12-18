package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/apis/v1/client"
	"github.com/sonujose/sloop/apis/v1/controller"
	"github.com/sonujose/sloop/pkg/core/consts"
	"github.com/sonujose/sloop/pkg/core/history"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
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

func (s *SyncConfig) SyncPackage(l *logrus.Logger) error {

	kcSecretBlob := make(map[string][]byte)

	deployedOn := time.Now()
	var lastSyncRevision int

	lastSyncRevision, err := getLastSyncRevision(s.kubeClient, s.sloopPkg.Metadata.Name, "sloop")
	if err != nil {
		l.Debugf("Unable to find any sync revison for the specified sloop package, creating the first revision. reason=%v", err)
	}

	newSyncRevision := lastSyncRevision + 1

	configStatus := controller.Status{
		DeployedOn:   fmt.Sprint(deployedOn),
		Version:      s.sloopPkg.Spec.Version,
		SyncRevision: newSyncRevision,
		Components:   len(s.sloopPkg.Spec.Components),
		Name:         s.sloopPkg.Metadata.Name,
	}

	s.controllerCfg.Status = configStatus
	blob, _ := json.Marshal(s.controllerCfg)
	kcSecretBlob["config"] = blob

	secretName := getSloopConfigSecretName(s.sloopPkg.Spec.Version, s.sloopPkg.Metadata.Name, newSyncRevision)

	s1 := &consts.SloopConfigSecretLabels{
		Components: fmt.Sprint(len(s.sloopPkg.Spec.Components)),
		Package:    s.sloopPkg.Metadata.Name,
		Revision:   fmt.Sprint(newSyncRevision),
		Status:     consts.StatusRegistered,
		Updated:    fmt.Sprint(deployedOn.Unix()),
		Version:    s.sloopPkg.Spec.Version,
	}

	SloopControllerSecret := &corev1.Secret{
		Data: kcSecretBlob,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: s.sloopPkg.Metadata.Namespace,
			Labels:    s1.GetConfigLabels(),
		},
		TypeMeta: metav1.TypeMeta{},
		Type:     corev1.SecretType(fmt.Sprintf("sloop.io/%s", s.sloopPkg.Metadata.Name)),
	}

	ctx := context.Background()
	_, err = s.kubeClient.CoreV1().Secrets(s.sloopPkg.Metadata.Namespace).Create(ctx, SloopControllerSecret, metav1.CreateOptions{})

	if err != nil {
		l.Errorf("Error creating sloop config secret.")
		return err
	}

	sloopDeploymentStatus, _ := yaml.Marshal(s.controllerCfg.Status)

	// CONSOLE-INFO : SYNC OPERATION STATUS
	fmt.Println("All done! Registered sloop configuration for controller")
	fmt.Println(string(sloopDeploymentStatus))

	return nil
}

func getLastSyncRevision(KubeClient *kubernetes.Clientset, name string, namespace string) (int, error) {

	hs := history.New(KubeClient)
	syncHistory, err := hs.GetPackageSyncHistory(name, namespace, false)

	var lastrev int = 0

	if err != nil {
		return lastrev, fmt.Errorf("Unable to fetch sync history from the cluster. error=%v", err)
	}

	var lastSyncExecutionMeta v1.Secret

	if len(syncHistory) == 0 {
		return lastrev, fmt.Errorf("No previous sync revisions found")
	}

	lastSyncExecutionMeta = syncHistory[0]

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
