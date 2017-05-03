package model

import (
	"log"
	"time"

	"path"

	"github.com/webdevwilson/terraform-ci/config"
	"github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/persist"
)

const projectNS = "projects"

// Project top-level data structure
type Project struct {
	GUID      string            `json:"guid,omitempty"`
	Name      string            `json:"name,omitempty"`
	RepoURL   string            `json:"repoUrl,omitempty"`
	RepoPath  string            `json:"repoPath,omitempty"`
	Settings  map[string]string `json:"settings,omitempty"`
	LocalPath string            `json:"-"`
}

var store persist.Store
var executor *execute.Executor

func init() {
	store = config.Get().Store
	store.CreateNamespace(projectNS)
	executor = config.Get().Executor

	ticker := time.NewTicker(time.Minute * 5)
	log.Printf("[INFO] Running plans every %d minutes", 5)
	go func() {
		for _ = range ticker.C {
			planRunner()
		}
	}()
	go planRunner()
}

// ListProjects returns all the projects
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

// GetProject fetches a project by guid
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

// CreateProject creates a new project
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

// ExecutePlan
func ExecutePlan(guid string) (taskID string, err error) {
	task := execute.Task{
		Command: "terraform",
		Args: []string{
			"apply",
			"-outfile",
			"terraform.tfplan",
		},
	}
	st, err := executor.Schedule(task)
	taskID = st.GUID
	return
}

// planRunner is responsible for starting up the goroutine that runs plans
func planRunner() {
	// Start ticker that runs plans
	prjs, err := ListProjects()
	if err != nil {
		log.Printf("[ERROR] Error retrieving projects: %s", err)
		return
	}

	// schedule job for each plan
	for _, prj := range prjs {
		go plan(&prj)
	}
}

func plan(prj *Project) {
	log.Printf("[INFO] Running plan for project '%s'", prj.GUID)
	task := execute.Task{
		Command: "terraform",
		Args: []string{
			"plan",
			"-out",
			path.Join(prj.LocalPath, "terraform.tfplan"),
			prj.LocalPath,
		},
	}
	st, err := executor.Schedule(task)
	if err != nil {
		log.Printf("[ERROR] Error scheduling plan run: %s", err)
		return
	}

	for result := range st.Channel {
		log.Printf("[INFO] Plan exited with code %d", result.ExitCode)
	}
}
