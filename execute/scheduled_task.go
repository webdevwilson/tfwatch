package execute

import (
	"fmt"
	"strings"
)

// ScheduledTask defines a task that has been scheduled to be executed on the system
type ScheduledTask struct {
	GUID string
	Task
	Channel      <-chan *Result
	writeChannel chan<- *Result
}

// String returns a string representation of the ScheduledTask
func (t ScheduledTask) String() string {
	return fmt.Sprintf("Task %s: '%s %s'", t.GUID, t.Command, strings.Join(t.Args, " "))
}

// GetTask returns a task from this
func (t ScheduledTask) task() Task {
	return Task{
		Command:          t.Command,
		Args:             t.Args,
		WorkingDirectory: t.WorkingDirectory,
		Environment:      t.Environment,
	}
}

// Result creates a result from this task
func (t ScheduledTask) Result(exitCode int, output []byte) *Result {
	return &Result{
		t.GUID,
		t.task(),
		exitCode,
		output,
	}
}
