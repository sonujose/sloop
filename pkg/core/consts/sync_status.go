package consts

const (
	// Sloop config status - registered : First phase of sloop lifecycle
	StatusRegistered string = "registered"
	// Sloop config status - synced : Successfully completed all components deployment
	StatusSynced string = "synced"
	// Sloop config status - syncError : Error during sync (no components deployed)
	StatusSyncError string = "syncError"
	// Sloop config Status - aborted : If the controller fails to pick up the revision before next new revision
	StatusAborted string = "aborted"
	// Sloop config status - partialSync : Only partial sync success (few components deployed)
	StatusCompleted string = "partialSync"
)
