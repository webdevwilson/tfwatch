package plan

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform/terraform"
	"github.com/webdevwilson/terraform-ui/config"
)

var planPath string

func init() {
	var err error
	planPath, err = config.Get().MakeStatePathDir("plans")
	if err != nil {
		log.Printf("[ERROR] Error creating path %s: %s", planPath, err)
	}
}

// List returns the name of plans waiting to be applied
func List() (plans []string, err error) {
	plans, err = filepath.Glob(path.Join(planPath, "*"))

	// remove the leading path
	for i, v := range plans {
		plans[i] = path.Base(v)
	}

	return
}

// Get returns a datastructure containing the plan
func Get(name string) (plan *terraform.Plan, err error) {
	file, err := os.Open(path.Join(planPath, name))

	if err != nil {
		return
	}

	plan, err = terraform.ReadPlan(file)

	return
}

// Apply applies a plan
func Apply(name string) (execID string, err error) {
	p := path.Join(planPath, name)
	file, err := os.Open(p)

	if err != nil {
		return
	}

	plan, err := terraform.ReadPlan(file)

	if err != nil {
		return
	}

	if plan.Diff.Empty() {
		return
	}

	task, err := executor.Schedule("terraform", "apply", p)

	if err != nil {
		log.Printf("[ERROR] Error scheduling task: %s", err)
		return
	}

	execID = task.GUID
	return
}
