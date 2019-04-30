package mkapi

// SyncTaskResults contains the results of a sync task.
type SyncTaskResults interface {
	// Good returns whether the task succeeded.
	Good() bool

	// Logs returns the task logs.
	Logs() string
}

// SyncTask is a generic sync task.
type SyncTask interface {
	// SetTimeout sets the task timeout in seconds.
	SetTimeout(timeoutSeconds int64)
}
