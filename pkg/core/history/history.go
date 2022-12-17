package history

import (
	"context"
	"sort"
	"strconv"

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

type SloopConfigSecretMeta struct {
	revision   int
	owner      string
	name       string
	modifiedAt int64
}

func (h *SloopHistory) GetPackageSyncHistory(name string, namespace string) ([]v1.Secret, error) {

	ctx := context.Background()

	lsel := kblabels.Set{"owner": "sloop", "name": name}.AsSelector()
	opts := metav1.ListOptions{LabelSelector: lsel.String()}

	secretList, err := h.KubeClient.CoreV1().Secrets(namespace).List(ctx, opts)

	if err != nil {
		return nil, err
	}

	sort.Slice(secretList.Items, func(i, j int) bool {
		l, _ := strconv.Atoi(secretList.Items[i].Labels["modifiedAt"])
		v, _ := strconv.Atoi(secretList.Items[j].Labels["modifiedAt"])
		return l > v
	})

	var secretMetaList = secretList.Items

	return secretMetaList, nil

}
