package task

import (
	"log"

	"os/exec"

	"syscall"

	"runtime/debug"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/webdevwilson/terraform-ci/persist"
)

const persistNamespace = "executions"

// Result contains the results of a task
type Result struct {
	ScheduledTask *ScheduledTask
	ExitCode      int
	Output        []byte
}

type Executor struct {
	store    persist.Store
	taskCh   chan ScheduledTask
	resultCh chan *Result
}

// NewExecutor creates returns a pointer to
func NewExecutor(store persist.Store) (exe *Executor) {
	taskQueue := make(chan ScheduledTask, 50)
	resultQueue := make(chan *Result, 50)

	exe = &Executor{
		store:    store,
		taskCh:   taskQueue,
		resultCh: resultQueue,
	}

	store.CreateNamespace(persistNamespace)

	go exe.runTasks()

	go exe.persistResults()

	return
}

// runTasks dequeues
func (exe *Executor) runTasks() {

	// gracefully recover from panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Executor recovered from panic: %s\n%s", r, debug.Stack())
			exe.runTasks()
		}
	}()

	for t := range exe.taskCh {
		log.Printf("[INFO] Executing %s", t.String())
		cmd := exec.Command(t.Command, t.Args...)

		output, err := cmd.CombinedOutput()

		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				log.Printf("[ERROR] Error executing %s: %s", t.String(), err)
			}
		}

		log.Printf("[INFO] Stdout: %s", output)

		// read the exit code
		_, ok := cmd.ProcessState.Sys().(*syscall.WaitStatus)
		if !ok {
			log.Printf("[ERROR] Error reading process status.")
		}

		// result := &Result{
		// 	&t,
		// 	status.ExitStatus(),
		// 	output,
		// }
		// result.Task.writeChannel <- result
		// resultQueue <- result
	}
}

// persistResults reads from resultChannel
func (exe *Executor) persistResults() {
	for r := range exe.resultCh {
		exe.store.Create(persistNamespace, r)
	}
}

// Schedule schedules a job to be run
func (exe *Executor) Schedule(command string, args ...string) (st *ScheduledTask, err error) {
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

	log.Printf("[INFO] Scheduling %s", st.String())

	// exe.taskCh <- st

	return
}

// GetResult returns a result from disk
func (exe *Executor) GetResult(guid string) (r Result, err error) {
	// f, err := os.Open(path.Join(resultPath, guid))

	if err != nil {
		return
	}

	// read in the value
	// err = gob.NewDecoder(f).Decode(&r)

	return
}
