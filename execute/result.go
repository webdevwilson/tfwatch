package execute

// Result contains the results of a task
type Result struct {
	GUID string
	Task
	ExitCode int
	Output   []byte
}
