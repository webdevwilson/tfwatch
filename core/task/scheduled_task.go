package task

import (
	"fmt"
	"strings"
)

// ScheduledTask defines a task that will be executed on the system
type ScheduledTask struct {
	GUID         string
	Command      string
	Args         []string
	Channel      <-chan *Result
	writeChannel chan<- *Result
}

// String returns a string representation of the ScheduledTask
func (t ScheduledTask) String() string {
	return fmt.Sprintf("Task %s: '%s %s'", t.GUID, t.Command, strings.Join(t.Args, " "))
}
