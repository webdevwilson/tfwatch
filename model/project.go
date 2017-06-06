package model

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"path"

	"github.com/hashicorp/terraform/terraform"
)

type ProjectStatus string

const (
	ProjectStatusNew     ProjectStatus = "new"
	ProjectStatusError   ProjectStatus = "error"
	ProjectStatusOK      ProjectStatus = "ok"
	ProjectStatusPending ProjectStatus = "pending"
)

// Project top-level data structure
type Project struct {
	GUID           string            `json:"guid,omitempty"`
	Name           string            `json:"name,omitempty"`
	Settings       map[string]string `json:"settings,omitempty"`
	PlanUpdated    time.Time         `json:"plan_updated,omitempty"`
	PendingChanges []ResourceChange  `json:"pending_changes"`
	Status         ProjectStatus     `json:"status,omitempty"`
	LocalPath      string            `json:"-"`
}

// ResourceChange represents a change
type ResourceChange struct {
	ResourceID string `json:"resource_id"`
	Action     string `json:"action"`
}

// NewProject creates a new project
func NewProject(name string, dir string) *Project {
	return &Project{
		Name:           name,
		LocalPath:      dir,
		Status:         "new",
		PendingChanges: []ResourceChange{},
	}
}

// Plan returns the Terraform plan for a project
func (prj *Project) Plan() (*Plan, error) {
	planFile := path.Join(prj.LocalPath, "terraform.tfplan")

	// Open the path no matter if its a directory or file
	f, err := os.Open(planFile)
	defer func() {
		log.Printf("[DEBUG] Closing plan file '%s'", planFile)
		f.Close()
	}()

	if err != nil {
		return nil, fmt.Errorf(
			"Failed to load Terraform configuration or plan: %s", err)
	}

	// Stat it so we can check if its a directory
	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to load Terraform configuration or plan: %s", err)
	}

	// If this path is a directory, then it can't be a plan. Not an error.
	if fi.IsDir() {
		return nil, fmt.Errorf(
			"Failed to load plan '%s', expected plan file, found directory", planFile)
	}

	// Read the plan
	log.Printf("[DEBUG] Reading plan from file '%s'", planFile)
	p, err := terraform.ReadPlan(f)
	log.Printf("[DEBUG] Plan read from file '%s'", planFile)
	if err != nil {
		return nil, err
	}

	return &Plan{p}, nil
}

// ExecutionNS returns the namespace to use for this projects executions
func (prj *Project) ExecutionNS() string {
	return fmt.Sprintf("project-%s-executions", prj.GUID)
}

// Plan is used to wrap a terraform.Plan and add methods
type Plan struct {
	plan *terraform.Plan
}

// ResourceChanges returns the changes in a plan
func (p *Plan) ResourceChanges() []*ResourceChange {
	var changes []*ResourceChange
	for _, module := range p.plan.Diff.Modules {
		for id, res := range module.Resources {
			if change := res.ChangeType(); change != terraform.DiffNone {
				var changeStr string
				switch change {
				case terraform.DiffCreate:
					changeStr = "Create"
				case terraform.DiffDestroyCreate:
					changeStr = "Recreate"
				case terraform.DiffDestroy:
					changeStr = "Destroy"
				case terraform.DiffUpdate:
					changeStr = "Update"
				}
				id := fmt.Sprintf("%s.%s", strings.Join(module.Path, "."), id)
				changes = append(changes, &ResourceChange{
					ResourceID: id,
					Action:     changeStr,
				})
			}
		}
	}
	return changes
}
