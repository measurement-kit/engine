// +build android

package engine

import "github.com/measurement-kit/engine/task"

// An asynchronous task handle.
type MKEAsyncTask struct {
	at *task.Handle
}

// An engine for running async tasks.
type MKEAsyncEngine struct {
}

// Starts a task configured according to settings.
func (ae *MKEAsyncEngine) Start(settings string) *MKEAsyncTask {
	return &MKEAsyncTask{at: task.Start(settings)}
}

// Indicates whether the task has finished running.
func (at *MKEAsyncTask) IsDone() bool {
	return at.at.IsDone()
}

// Returns the next event emitted by the task.
func (at *MKEAsyncTask) WaitForNextEvent() string {
	return at.at.WaitForNextEvent()
}

// Attempts to interrupt a running task.
func (at *MKEAsyncTask) Interrupt() {
	at.at.Interrupt()
}
