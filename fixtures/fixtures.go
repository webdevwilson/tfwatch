package fixtures

import (
	"os"
	"path"

	"github.com/webdevwilson/tfwatch/model"
)

type Fixture string

const TerraformNoPlan Fixture = "terraform_noplan"
const TerraformPlanned Fixture = "terraform_planned"

// Returns a project from the filesystem
func GetProject(name Fixture) (*model.Project, error) {
	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return model.NewProject(string(name), path.Join(wd, string(name))), nil
}
