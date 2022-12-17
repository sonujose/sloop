package history

import (
	"context"
	"sort"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kblabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type SloopHistory struct {
	KubeClient *kubernetes.Clientset
}

func New(kclient *kubernetes.Clientset) *SloopHistory {
	return &SloopHistory{
		KubeClient: kclient,
	}
}

type SyncHistoryObj struct {
	Revision   string
	Package    string
	Updated    time.Time
	Version    string
	Components string
}

func (h *SloopHistory) ListSyncHistory(name string, namespace string, allpackages bool) ([]SyncHistoryObj, error) {

	var syncHistory []SyncHistoryObj
	syncSecretList, err := h.GetPackageSyncHistory(name, namespace, allpackages)

	if err != nil {
		return nil, err
	}

	for _, syncSecret := range syncSecretList {

		updatedOn := time.Unix(syncSecret.CreationTimestamp.Unix(), 0)

		sh := SyncHistoryObj{
			Revision:   syncSecret.Labels["revision"],
			Updated:    updatedOn,
			Package:    syncSecret.Labels["package"],
			Version:    syncSecret.Labels["version"],
			Components: syncSecret.Labels["components"],
		}
		syncHistory = append(syncHistory, sh)
	}

	return syncHistory, nil
}

func (h *SloopHistory) GetPackageSyncHistory(name string, namespace string, allpackages bool) ([]v1.Secret, error) {

	ctx := context.Background()

	var lsel kblabels.Selector
	if allpackages {
		lsel = kblabels.Set{"owner": "sloop"}.AsSelector()
	} else {
		lsel = kblabels.Set{"owner": "sloop", "package": name}.AsSelector()
	}

	opts := metav1.ListOptions{LabelSelector: lsel.String()}

	secretList, err := h.KubeClient.CoreV1().Secrets(namespace).List(ctx, opts)

	if err != nil {
		return nil, err
	}

	if len(secretList.Items) > 1 {
		sort.Slice(secretList.Items, func(i, j int) bool {
			l, _ := strconv.Atoi(secretList.Items[i].Labels["updated"])
			v, _ := strconv.Atoi(secretList.Items[j].Labels["updated"])
			return l > v
		})
	}

	var secretMetaList = secretList.Items

	return secretMetaList, nil

}
