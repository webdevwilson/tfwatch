package model

import (
	"github.com/webdevwilson/terraform-ci/config"
	"github.com/webdevwilson/terraform-ci/persist"
)

const projectNS = "projects"

// Project top-level data structure
type Project struct {
	GUID     string `json:"guid,omitempty"`
	Name     string `json:"name,omitempty"`
	RepoURL  string `json:"repoUrl,omitempty"`
	RepoPath string `json:"repoPath,omitempty"`
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
		projects[i].GUID = guid
		if err != nil {
			return
		}
	}

	return
}

// GetProject
func GetProject(guid string) (*Project, error) {
	var prj Project
	err := store.Get(projectNS, guid, &prj)

	if err != nil {
		return nil, err
	}

	prj.GUID = guid
	return &prj, err
}

// GetProjectByName returns the named project
func GetProjectByName(name string) (*Project, error) {

	projects, err := store.List(projectNS)
	if err != nil {
		return nil, err
	}

	for _, guid := range projects {
		prj, err := GetProject(guid)

		if err != nil {
			return nil, err
		}

		if prj.Name == name {
			return prj, nil
		}
	}
	return nil, nil
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
