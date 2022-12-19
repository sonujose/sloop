package history

import (
	"context"
	"sort"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klabel "k8s.io/apimachinery/pkg/labels"
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
	Status     string
}

func (h *SloopHistory) ListSyncHistory(name string, namespace string, allpackages bool) ([]SyncHistoryObj, error) {

	var syncHistory []SyncHistoryObj

	var labelSel klabel.Selector

	if allpackages {
		labelSel = GetPackageHistoryLabelSelectors(AllPackagesHistoryFilterKey, "")
	} else {
		labelSel = GetPackageHistoryLabelSelectors(PackageHistoryFilterKey, name)
	}

	pkgHistoryList, err := h.GetPackageSyncHistorybyLabels(labelSel, namespace)

	if err != nil {
		return nil, err
	}

	for _, syncSecret := range pkgHistoryList {

		updatedOn := time.Unix(syncSecret.CreationTimestamp.Unix(), 0)

		sh := SyncHistoryObj{
			Revision:   syncSecret.Labels["revision"],
			Updated:    updatedOn,
			Package:    syncSecret.Labels["package"],
			Version:    syncSecret.Labels["version"],
			Components: syncSecret.Labels["components"],
			Status:     syncSecret.Labels["status"],
		}
		syncHistory = append(syncHistory, sh)
	}

	return syncHistory, nil
}

// GetPackageSyncHistory - returnd the package history info in descending order of deployment
func (h *SloopHistory) GetPackageSyncHistorybyLabels(labelSelector klabel.Selector, namespace string) ([]v1.Secret, error) {

	ctx := context.Background()

	secretList, err := h.KubeClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector.String()})

	if err != nil {
		return nil, err
	}

	if len(secretList.Items) > 1 {
		sort.Slice(secretList.Items, func(i, j int) bool {
			l, _ := strconv.Atoi(secretList.Items[i].Labels["updated"])
			v, _ := strconv.Atoi(secretList.Items[j].Labels["updated"])
			return l < v
		})
	}

	var secretMetaList = secretList.Items

	return secretMetaList, nil

}
