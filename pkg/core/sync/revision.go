package sync

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sonujose/sloop/pkg/core/consts"
	"github.com/sonujose/sloop/pkg/core/history"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getLastSyncRevision - Returns the last sync version of the package
func getLastSyncRevision(syncHistory []v1.Secret) (int, error) {

	var lastrev int = 0
	var lastSyncExecutionMeta v1.Secret

	if len(syncHistory) == 0 {
		return lastrev, fmt.Errorf("No previous sync revisions found")
	}

	lastSyncExecutionMeta = syncHistory[len(syncHistory)-1]

	lastrev, err := strconv.Atoi(lastSyncExecutionMeta.Labels["revision"])
	if err != nil {
		return lastrev, err
	}

	return lastrev, nil
}

// cleanOldRegisteredRevisions - cleans up old revisions (making the status to aborted), which are not yet catched by the controller, so that
// only the latest rev will be picked up by the controller.
func (s *SyncConfig) cleanOldRegisteredRevisions(rev string) error {
	hs := history.New(s.kubeClient, s.log)

	var syncHistoryPending []v1.Secret

	labelSel := history.GetPackageHistoryLabelSelectors(history.RegisteredPackagesHistoryFilterKey, s.sloopPkg.Metadata.Name)
	syncHistoryPending, err := hs.GetPackageSyncHistorybyLabels(labelSel, s.sloopPkg.Metadata.Namespace)

	if err != nil {
		return err
	}

	for _, secretCopy := range syncHistoryPending {
		if secretCopy.Labels[consts.ConfigLabelRevision] != rev {

			secretCopy.Labels[consts.ConfigLabelStatus] = consts.StatusAborted

			// Update secret with new label status
			_, err := s.kubeClient.CoreV1().Secrets(s.sloopPkg.Metadata.Namespace).Update(context.Background(), &secretCopy, metav1.UpdateOptions{})
			if err != nil {
				s.log.WithError(err).Debugf("Unable to update old secret revisions")
				continue
			}
		}

	}

	return nil

}
