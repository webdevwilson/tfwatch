package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"path"

	"github.com/hashicorp/terraform/terraform"
)

// Project top-level data structure
type Project struct {
	GUID        string            `json:"guid,omitempty"`
	Name        string            `json:"name,omitempty"`
	RepoURL     string            `json:"repo_url,omitempty"`
	RepoPath    string            `json:"repo_path,omitempty"`
	Settings    map[string]string `json:"settings,omitempty"`
	PlanUpdated time.Time         `json:"plan_updated,omitempty"`
	Status      string            `json:"status,omitempty"`
	LocalPath   string            `json:"-"`
}

// NewProject creates a new project
func NewProject(name string, dir string) *Project {
	return &Project{
		Name:      name,
		LocalPath: dir,
		Status:    "new",
	}
}

// Plan returns the Terraform plan for a project
func (prj *Project) Plan() (*terraform.Plan, error) {
	planFile := path.Join(prj.LocalPath, "terraform.tfplan")
	plan, err := os.Open(planFile)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Reading plan from file '%s'", planFile)
	return terraform.ReadPlan(plan)
}

// ExecutionNS returns the namespace to use for this projects executions
func (prj *Project) ExecutionNS() string {
	return fmt.Sprintf("project-%s-executions", prj.GUID)
}
