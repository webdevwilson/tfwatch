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

// Executor runs processes on the machine and persists results
type Executor interface {
	Schedule(Task) (*ScheduledTask, error)
}

// Executor is used to schedule tasks to run
type executor struct {
	store    persist.Store
	logDir   string
	taskCh   chan *ScheduledTask
	resultCh chan *Result
}

// NewExecutor creates returns a pointer to
func NewExecutor(store persist.Store, logDir string) Executor {

	log.Printf("[INFO] Executor log directory: %s", logDir)

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Printf("[ERROR] Error creating executor log directory %s: %s", logDir, err)
	}

	exe := &executor{
		store:    store,
		logDir:   logDir,
		taskCh:   make(chan *ScheduledTask, 50),
		resultCh: make(chan *Result, 50),
	}

	store.CreateNamespace(persistNamespace)

	go exe.runTasks()

	go exe.persistResults()

	return exe
}

// runTasks dequeues
func (exe *executor) runTasks() {

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

		// in cases where the command was not executed (not found on path)
		// exit code is -1 and output is the error message
		var statusCode int
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				log.Printf("[ERROR] Error executing %s: %s", t.String(), err)
				statusCode = -1
				output = []byte(err.Error())
				return
			}
		}

		// log output to disk
		logFile := path.Join(exe.logDir, t.GUID)
		err = ioutil.WriteFile(logFile, output, os.ModePerm)
		if err != nil {
			log.Printf("[WARN] Error logging task output to disk: %s", err)
		}

		// read the exit code
		if statusCode == 0 {
			status, ok := cmd.ProcessState.Sys().(*syscall.WaitStatus)
			if !ok {
				log.Printf("[ERROR] Error reading process status for task '%s'", t.GUID)
				statusCode = -2
			} else {
				statusCode = status.ExitStatus()
			}
		}

		// create result, and send across channel
		result := t.Result(statusCode, output)
		t.writeChannel <- result
		exe.resultCh <- result
	}
}

// persistResults reads from resultChannel
func (exe *executor) persistResults() {
	for r := range exe.resultCh {

		log.Printf("[INFO] Persisting result for task '%s'", r.GUID)

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
func (exe *executor) Schedule(task Task) (st *ScheduledTask, err error) {

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
func (exe *executor) GetResult(guid string) (r Result, err error) {
	// f, err := os.Open(path.Join(resultPath, guid))

	if err != nil {
		return
	}

	// read in the value
	// err = gob.NewDecoder(f).Decode(&r)

	return
}
