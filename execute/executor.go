package execute

import (
	"log"
	"path"

	"os/exec"

	"syscall"

	"runtime/debug"

	"os"

	"io/ioutil"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/webdevwilson/terraform-ci/persist"
)

const persistNamespace = "executions"

// Executor is used to schedule tasks to run
type Executor struct {
	store    persist.Store
	logDir   string
	taskCh   chan *ScheduledTask
	resultCh chan *Result
}

// NewExecutor creates returns a pointer to
func NewExecutor(store persist.Store, logDir string) (exe *Executor) {

	log.Printf("[INFO] Executor log directory: %s", logDir)

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Printf("[ERROR] Error creating executor log directory %s: %s", logDir, err)
	}

	exe = &Executor{
		store:    store,
		logDir:   logDir,
		taskCh:   make(chan *ScheduledTask, 50),
		resultCh: make(chan *Result, 50),
	}

	store.CreateNamespace(persistNamespace)

	go exe.runTasks()

	go exe.persistResults()

	return
}

// runTasks dequeues
func (exe *Executor) runTasks() {

	// gracefully recover from panics and continue to run tasks
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Executor recovered from panic: %s\n%s", r, debug.Stack())
			exe.runTasks()
		}
	}()

	for t := range exe.taskCh {
		log.Printf("[INFO] Executing %s", t.String())
		cmd := exec.Command(t.Command, t.Args...)

		if t.WorkingDirectory != "" {
			cmd.Dir = t.WorkingDirectory
		}

		output, err := cmd.CombinedOutput()

		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				log.Printf("[ERROR] Error executing %s: %s", t.String(), err)
			}
		}

		// log output to disk
		logFile := path.Join(exe.logDir, t.GUID)
		err = ioutil.WriteFile(logFile, output, os.ModePerm)
		if err != nil {
			log.Printf("[WARN] Error logging task output to disk: %s", err)
		}

		// read the exit code
		var statusCode int
		status, ok := cmd.ProcessState.Sys().(*syscall.WaitStatus)
		if !ok {
			log.Printf("[ERROR] Error reading process status.")
			statusCode = -99
		} else {
			statusCode = status.ExitStatus()
		}

		result := t.Result(statusCode, output)

		t.writeChannel <- result
		exe.resultCh <- result
	}
}

// persistResults reads from resultChannel
func (exe *Executor) persistResults() {
	for r := range exe.resultCh {

		log.Printf("[INFO] Persisting result for task %s", r.GUID)

		// persist result
		exe.store.Create(persistNamespace, r)

		// output logs
		logFile := path.Join(exe.logDir, r.GUID)
		err := ioutil.WriteFile(logFile, r.Output, os.ModePerm)
		if err != nil {
			log.Printf("[WARN] Failed to write to log file '%s': %s", logFile, err)
		}
	}
}

// Schedule schedules a job to be run
func (exe *Executor) Schedule(task Task) (st *ScheduledTask, err error) {

	// Create a GUID for the task
	uidPtr, err := uuid.NewV4()
	if err != nil {
		return
	}

	// Read / Write channel is same
	ch := make(chan *Result, 1)
	st = &ScheduledTask{
		uidPtr.String(),
		task,
		ch,
		ch,
	}

	log.Printf("[INFO] Scheduling %s", st.String())

	exe.taskCh <- st

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
