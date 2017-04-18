package model

import (
	"github.com/webdevwilson/terraform-ui/config"
	"github.com/webdevwilson/terraform-ui/persist"
)

const projectNS = "projects"

type Project struct {
	GUID     string `json:"omitempty"`
	Name     string
	RepoURL  string
	RepoPath string
}

var store persist.Store

func init() {
	store = config.Get().Store
	store.CreateNamespace(projectNS)
}

// ListProjects
func ListProjects() (projects []Project, err error) {
	guids, err := store.List(projectNS)

	if err != nil {
		return
	}

	projects = make([]Project, len(guids))
	for i, guid := range guids {
		err = store.Get(projectNS, guid, &projects[i])
		if err != nil {
			return
		}
	}

	return
}

// GetProject
func GetProject(guid string) (*Project, error) {
	var prj *Project
	err := store.Get(projectNS, guid, prj)

	if err != nil {
		return nil, err
	}

	prj.GUID = guid
	return prj, err
}

// CreateProject
func CreateProject(prj *Project) (err error) {
	var guid string
	guid, err = store.Create(projectNS, prj)

	if err != nil {
		return
	}

	prj.GUID = guid
	return
}

// UpdateProject
func UpdateProject(prj *Project) error {
	return store.Update(projectNS, prj.GUID, prj)
}

// DeleteProject
func DeleteProject(guid string) error {
	return store.Delete(projectNS, guid)
}