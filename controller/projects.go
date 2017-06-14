package controller

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"time"

	"github.com/webdevwilson/tfwatch/execute"
	"github.com/webdevwilson/tfwatch/model"
	"github.com/webdevwilson/tfwatch/persist"
)

const projectNS = "projects"

var bootstrapIgnores = []string{".git", ".tfwatch", "node_modules"}

// Projects hosts the business logic for projects
type Projects interface {
	List() (projects []*model.Project, err error)
	Get(guid string) (*model.Project, error)
	GetByName(name string) (*model.Project, error)
	Create(prj *model.Project) (err error)
	Update(prj *model.Project) error
	Delete(guid string) error
	ExecutePlan(prj *model.Project) (taskID string, err error)
	GetExecutions(prj *model.Project) (results []*execute.Result, err error)
}

type projects struct {
	store        persist.Store
	executor     execute.Executor
	planInterval time.Duration
	runPlans     bool
}

// NewProjectsController creates a new controller
func NewProjectsController(dir string, store persist.Store, executor execute.Executor,
	interval time.Duration, runPlans bool) Projects {

	store.CreateNamespace(projectNS)

	p := &projects{
		store:        store,
		executor:     executor,
		planInterval: interval,
		runPlans:     runPlans,
	}

	// Start plans for existing projects
	prjs, err := p.List()
	if err != nil {
		log.Printf("[ERROR] Error scheduling project plan: %s", err)
	}

	for _, prj := range prjs {
		p.schedulePlan(interval, prj)
	}

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
			p.Create(model.NewProject(name, dir))
		}
	}
}

// ListProjects returns all the projects
func (p *projects) List() (projects []*model.Project, err error) {
	guids, err := p.store.List(projectNS)

	if err != nil {
		return
	}

	projects = make([]*model.Project, len(guids))
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

	// create namespace to store executions for the project
	err = p.store.CreateNamespace(prj.ExecutionNS())
	if err != nil {
		return
	}

	// schedule plan updates
	p.schedulePlan(p.planInterval, prj)

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
func (p *projects) ExecutePlan(prj *model.Project) (string, error) {
	taskID, _, err := p.executeInProject(prj, &execute.Task{
		Command: "terraform",
		Args: []string{
			"apply",
			"-outfile",
			"terraform.tfplan",
		},
	})

	return taskID, err
}

// GetExecutions returns the executions that have occurred in a project
func (p *projects) GetExecutions(prj *model.Project) (r []*execute.Result, err error) {
	guids, err := p.store.List(prj.ExecutionNS())
	if err != nil {
		return
	}
	r = make([]*execute.Result, len(guids))

	for i, guid := range guids {
		p.store.Get(prj.ExecutionNS(), guid, &r[i])
	}

	return
}

// executeInProject
func (p *projects) executeInProject(prj *model.Project, t *execute.Task) (taskID string, ch <-chan *execute.Result, err error) {
	t.WorkingDirectory = prj.LocalPath
	st, err := p.executor.Schedule(t)
	taskID = st.GUID
	writeCh := make(chan *execute.Result, 1)
	ch = writeCh

	// persist execution
	go func() {
		r := <-st.Channel
		_, err := p.store.Create(prj.ExecutionNS(), r)
		if err != nil {
			log.Printf("[ERROR] Error persisting execution in project '%s': %s", prj.GUID, err)
		}
		writeCh <- r
	}()

	return
}
