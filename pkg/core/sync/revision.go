package sync

import (
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
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

func cleanOldRegisteredRevisions(syncHistory []v1.Secret) {

}
