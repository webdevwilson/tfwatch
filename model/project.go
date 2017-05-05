package model

import (
	"log"
	"os"

	"path"

	"github.com/hashicorp/terraform/terraform"
)

// Project top-level data structure
type Project struct {
	GUID      string            `json:"guid,omitempty"`
	Name      string            `json:"name,omitempty"`
	RepoURL   string            `json:"repoUrl,omitempty"`
	RepoPath  string            `json:"repoPath,omitempty"`
	Settings  map[string]string `json:"settings,omitempty"`
	LocalPath string            `json:"-"`
}

// returns the plan for a project
func (prj *Project) Plan() (*terraform.Plan, error) {
	planFile := path.Join(prj.LocalPath, "terraform.tfplan")
	plan, err := os.Open(planFile)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Reading plan from file '%s'", planFile)
	return terraform.ReadPlan(plan)
}
