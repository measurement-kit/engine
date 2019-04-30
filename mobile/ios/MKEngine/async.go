// +build ios

package MKEngine

import "github.com/measurement-kit/engine/task"

// An asynchronous task handle.
type AsyncTask struct {
	at *task.Handle
}

// An engine for running async tasks.
type AsyncEngine struct {
}

// Starts a task configured according to settings.
func (ae *AsyncEngine) Start(settings string) *AsyncTask {
	return &AsyncTask{at: task.Start(settings)}
}

// Indicates whether the task has finished running.
func (at *AsyncTask) IsDone() bool {
	return at.at.IsDone()
}

// Returns the next event emitted by the task.
func (at *AsyncTask) WaitForNextEvent() string {
	return at.at.WaitForNextEvent()
}

// Attempts to interrupt a running task.
func (at *AsyncTask) Interrupt() {
	at.at.Interrupt()
}
