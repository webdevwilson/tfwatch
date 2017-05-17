package execute

// Task
type Task struct {
	Command          string
	Args             []string
	WorkingDirectory string
	Environment      map[string]string
}
