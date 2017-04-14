package task

import (
	"log"
	"os"
	"path"

	"os/exec"

	"encoding/gob"

	"syscall"

	"runtime/debug"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/webdevwilson/terraform-ui/core/persist"
)

const persistNamespace = "executions"

// Executor is responsible for executing commands and persisting results
type Executor interface {
	Schedule(cmd string, args ...string) (*ScheduledTask, error)
	GetResult(guid string) (Result, error)
}

type executor struct {
	store    *persist.Store
	taskCh   chan Task
	resultCh chan *Result
}

// NewExecutor creates returns a pointer to
func NewExecutor(store persist.Store) (exe Executor) {
	taskQueue := make(chan Task, 50)
	resultQueue := make(chan *Result, 50)

	exe = &executor{
		store:       store,
		taskQueue:   taskQueue,
		resultQueue: resultQueue,
	}

	store.CreateNamespace(persistNamespace)

	go exe.runTasks()

	go exe.persistResults()

	return
}

// runTasks dequeues
func (exe executor) runTasks() {

	// gracefully recover from panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Executor recovered from panic: %s\n%s", r, debug.Stack())
			exe.runTasks()
		}
	}()

	for t := range exe.taskCh {
		log.Printf("[INFO] Executing %s", t.string())
		cmd := exec.Command(t.Command, t.Args...)

		output, err := cmd.CombinedOutput()

		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				log.Printf("[ERROR] Error executing %s: %s", t.string(), err)
			}
		}

		log.Printf("[INFO] Stdout: %s", output)

		// read the exit code
		status, ok := cmd.ProcessState.Sys().(*syscall.WaitStatus)
		if !ok {
			log.Printf("[ERROR] Error reading process status.")
		}

		result := &Result{
			&t,
			status.ExitStatus(),
			output,
		}
		result.Task.writeChannel <- result
		resultQueue <- result
	}
}

// persistResults reads from resultChannel
func (exe executor) persistResults() {
	for r := range exe.resultCh {
		exe.store.Save(persistNamespace, r)
	}
}

// Schedule schedules a job to be run
func (exe executor) Schedule(command string, args ...string) (st ScheduledTask, err error) {
	uidPtr, err := uuid.NewV4()

	if err != nil {
		return
	}

	ch := make(chan *Result, 1)
	st = &ScheduledTask{
		uidPtr.String(),
		command,
		args,
		ch,
		ch,
	}

	log.Printf("[INFO] Scheduling %s", t.string())

	taskQueue <- t

	return
}

// GetResult returns a result from disk
func (exe executor) GetResult(guid string) (r Result, err error) {
	f, err := os.Open(path.Join(resultPath, guid))

	if err != nil {
		return
	}

	// read in the value
	err = gob.NewDecoder(f).Decode(&r)

	return
}
