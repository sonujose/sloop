package sync

import (
	"fmt"
	"strings"

	"github.com/sonujose/sloop/apis/v1/controller"
	"github.com/sonujose/sloop/pkg/core/consts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
