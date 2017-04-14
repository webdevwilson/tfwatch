package project

import (
	"fmt"
	"log"
	"path"

	"github.com/webdevwilson/terraform-ui/config"
	"github.com/webdevwilson/terraform-ui/core/persist"
)

// Model is the data structure for persisting
type Model struct {
	GUID       string
	Name       string
	SourceURL  string
	SourcePath string
}

var store *persist.Store

const storeNamespace = "store"

func init() {
	settings := config.Get()
	store = settings.Store

	err = store.CreateNamespace(storeNamespace)
	if err != nil {
		log.Printf("[ERROR] Error creating store namespace '%s': %s", storeNamespace, err)
	}
}

// List returns the projects in the system, when error returned, nil is always returned
func List() (projects []*Model, err error) {

	guids, err = store.List(storeNamespace)
	if err != nil {
		return
	}

	projects = make([]*Model, len(projects))

	for i, guid := range guids {
		err = store.Get(storeNamespace, &projects[i])
		if err != nil {
			projects = nil
			return
		}
	}

	return
}

// Diff checks a project for diffs
func Diff(name string) (err error) {
	p := path.Join(projectPath, name)
	planFile := path.Join(p, fmt.Sprintf("%s.tfplan", name))
	_, err = executor.Schedule("terraform", "plan", "-out", planFile, p)
	return
}
