package history

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sonujose/sloop/pkg/core/consts"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klabel "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type History struct {
	KubeClient *kubernetes.Clientset
	log        *logrus.Logger
}

func New(kclient *kubernetes.Clientset, log *logrus.Logger) *History {
	return &History{
		KubeClient: kclient,
		log:        log,
	}
}

type SyncHistory struct {
	Revision   string
	Package    string
	Updated    time.Time
	Version    string
	Components string
	Status     string
}

func (h *History) ListSyncHistory(name string, namespace string, allpackages bool) ([]SyncHistory, error) {

	var syncHistory []SyncHistory

	var labelSel klabel.Selector

	if allpackages {
		labelSel = h.GetPackageHistoryLabelSelectors(AllPackagesHistoryFilterKey, "", "")
	} else {
		labelSel = h.GetPackageHistoryLabelSelectors(PackageHistoryFilterKey, name, "")
	}

	pkgHistoryList, err := h.GetPackageSyncHistorybyLabels(labelSel, namespace)

	if err != nil {
		return nil, err
	}

	for _, syncSecret := range pkgHistoryList {

		updatedOn := time.Unix(syncSecret.CreationTimestamp.Unix(), 0)

		sh := SyncHistory{
			Revision:   syncSecret.Labels[consts.ConfigLabelRevision],
			Updated:    updatedOn,
			Package:    syncSecret.Labels[consts.ConfigLabelPackage],
			Version:    syncSecret.Labels[consts.ConfigLabelVersion],
			Components: syncSecret.Labels[consts.ConfigLabelComponents],
			Status:     syncSecret.Labels[consts.ConfigLabelStatus],
		}
		syncHistory = append(syncHistory, sh)
	}

	return syncHistory, nil
}

// GetPackageSyncHistory - returnd the package history info in descending order of deployment
func (h *History) GetPackageSyncHistorybyLabels(labelSelector klabel.Selector, namespace string) ([]v1.Secret, error) {

	ctx := context.Background()

	secretList, err := h.KubeClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector.String()})

	if err != nil {
		return nil, err
	}

	if len(secretList.Items) > 1 {
		sort.Slice(secretList.Items, func(i, j int) bool {
			l, _ := strconv.Atoi(secretList.Items[i].Labels[consts.ConfigLabelUpdated])
			v, _ := strconv.Atoi(secretList.Items[j].Labels[consts.ConfigLabelUpdated])
			return l < v
		})
	}

	var secretMetaList = secretList.Items

	return secretMetaList, nil

}
