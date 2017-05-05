package controller

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/model"
	"github.com/webdevwilson/terraform-ci/persist"
)

const projectNS = "projects"

var bootstrapIgnores = []string{".git", ".terraform-ci", "node_modules"}

// Projects hosts the business logic for projects
type Projects interface {
	List() (projects []model.Project, err error)
	Get(guid string) (*model.Project, error)
	GetByName(name string) (*model.Project, error)
	Create(prj *model.Project) (err error)
	Update(prj *model.Project) error
	Delete(guid string) error
	ExecutePlan(guid string) (taskID string, err error)
}

type projects struct {
	store      persist.Store
	executor   execute.Executor
	planTicker *time.Ticker
}

// NewProjectsController creates a new controller
func NewProjectsController(dir string, store persist.Store, executor execute.Executor) Projects {

	store.CreateNamespace(projectNS)

	ticker := time.NewTicker(time.Minute * 5)
	p := &projects{
		store:      store,
		executor:   executor,
		planTicker: ticker,
	}

	log.Printf("[INFO] Running plans every %d minutes", 5)
	go func() {
		for _ = range ticker.C {
			p.planRunner()
		}
	}()
	go p.planRunner()

	go p.bootstrap(dir)

	return p
}

func (p *projects) bootstrap(dir string) {
	// Naive scan for terraform projects
	log.Printf("[DEBUG] Bootstrapping projects in %s", dir)

	dirs, err := filepath.Glob(path.Join(dir, "*"))
	if err != nil {
		log.Fatalf("[FATAL] Error encountered loading projects: %s", err)
		return
	}

	for _, dir := range dirs {

		// search ignores
		var ignore bool
		for i := range bootstrapIgnores {
			if path.Base(dir) == bootstrapIgnores[i] {
				ignore = true
				break
			}
		}
		if ignore {
			continue
		}

		// is this even a directory?
		fi, err := os.Stat(dir)
		if err != nil {
			log.Printf("[ERROR] Error determining if '%s' is directory", dir)
			continue
		}

		if !fi.IsDir() {
			continue
		}

		// does this path have .tf files?
		tf, err := filepath.Glob(path.Join(dir, "*.tf"))
		if err != nil {
			log.Printf("[ERROR] Error encountered searching for .tf files: %s", err)
			continue
		}

		// none found? continue to next
		if len(tf) == 0 {
			log.Printf("[DEBUG] No .tf files found in %s, skipping", dir)
			continue
		}

		name := path.Base(dir)
		prj, err := p.GetByName(name)

		if err != nil {
			log.Printf("[ERROR] Error bootstrapping projects: %s", err)
			continue
		}

		if prj == nil {
			log.Printf("[INFO] Creating project '%s' found at '%s'", name, dir)
			p.Create(&model.Project{
				Name:      name,
				LocalPath: dir,
			})
		}
	}
}

// ListProjects returns all the projects
func (p *projects) List() (projects []model.Project, err error) {
	guids, err := p.store.List(projectNS)

	if err != nil {
		return
	}

	projects = make([]model.Project, len(guids))
	for i, guid := range guids {
		err = p.store.Get(projectNS, guid, &projects[i])
		projects[i].GUID = guid
		if err != nil {
			return
		}
	}

	return
}

// GetProject fetches a project by guid
func (p *projects) Get(guid string) (*model.Project, error) {
	var prj model.Project
	err := p.store.Get(projectNS, guid, &prj)

	if err != nil {
		return nil, err
	}

	prj.GUID = guid
	return &prj, err
}

// GetProjectByName returns the named project
func (p *projects) GetByName(name string) (*model.Project, error) {

	projects, err := p.store.List(projectNS)
	if err != nil {
		return nil, err
	}

	for _, guid := range projects {
		prj, err := p.Get(guid)

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
func (p *projects) Create(prj *model.Project) (err error) {
	var guid string
	guid, err = p.store.Create(projectNS, prj)

	if err != nil {
		return
	}

	prj.GUID = guid

	return
}

// UpdateProject
func (p *projects) Update(prj *model.Project) error {
	return p.store.Update(projectNS, prj.GUID, prj)
}

// DeleteProject
func (p *projects) Delete(guid string) error {
	return p.store.Delete(projectNS, guid)
}

// ExecutePlan
func (p *projects) ExecutePlan(guid string) (taskID string, err error) {
	task := execute.Task{
		Command: "terraform",
		Args: []string{
			"apply",
			"-outfile",
			"terraform.tfplan",
		},
	}
	st, err := p.executor.Schedule(task)
	taskID = st.GUID
	return
}

// planRunner is responsible for starting up the goroutine that runs plans
func (p *projects) planRunner() {
	// Start ticker that runs plans
	prjs, err := p.List()
	if err != nil {
		log.Printf("[ERROR] Error retrieving projects: %s", err)
		return
	}

	// schedule job for each plan
	for _, prj := range prjs {
		go p.plan(&prj)
	}
}

func (p *projects) plan(prj *model.Project) {
	log.Printf("[INFO] Running plan for project '%s'", prj.GUID)
	task := execute.Task{
		Command:          "terraform",
		WorkingDirectory: prj.LocalPath,
		Args:             []string{"plan", "-out", "terraform.tfplan"},
	}
	st, err := p.executor.Schedule(task)
	if err != nil {
		log.Printf("[ERROR] Error scheduling plan run: %s", err)
		return
	}

	for result := range st.Channel {
		log.Printf("[INFO] Plan exited with code %d", result.ExitCode)
	}
}
