package task

// Result contains the results of a task
type Result struct {
	Execution *Execution
	ExitCode  int
	Output    []byte
}
