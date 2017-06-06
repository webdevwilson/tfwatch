package controller

import (
	"log"
	"time"

	"github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/model"
)

// schedulePlan schedules a plan to run in a project, then waits and schedules it again
func (p *projects) schedulePlan(interval time.Duration, prj *model.Project) {

	// don't schedule if it has been configured not to run
	if !p.runPlans {
		log.Printf("[DEBUG] Plan running is disabled, do not schedule")
		return
	}

	go func() {
		p.runPlan(prj)
		time.AfterFunc(interval, func() {
			<-p.runPlan(prj)
			p.schedulePlan(interval, prj)
		})
	}()
}

func (p *projects) runPlan(prj *model.Project) (done <-chan bool) {
	log.Printf("[INFO] Running plan for project '%s'", prj.GUID)
	task := &execute.Task{
		Command: "terraform",
		Args: []string{
			"plan",
			"-detailed-exitcode",
			"-out",
			"terraform.tfplan",
		},
	}
	_, ch, err := p.executeInProject(prj, task)
	if err != nil {
		log.Printf("[ERROR] Error scheduling plan run: %s", err)
		return
	}

	// when task is complete, update the project
	done = make(chan bool, 1)
	go p.planComplete(prj, ch, done)

	return
}

func (p *projects) planComplete(prj *model.Project, ch <-chan *execute.Result, doneCh <-chan bool) {

	// wait for result
	r := <-ch

	// update project
	prj.PlanUpdated = time.Now()
	switch r.ExitCode {
	case 0:
		prj.Status = model.ProjectStatusOK
	case 2:
		prj.Status = model.ProjectStatusPending
	default:
		prj.Status = model.ProjectStatusError
		log.Printf("[WARN] Plan failed on %s: %s", prj.Name, r.Output)
	}
	log.Printf("[INFO] Project '%s' plan complete, updating status to '%s'", prj.GUID, prj.Status)

	// read the plan and add changes to project
	if prj.Status == model.ProjectStatusPending {
		plan, err := prj.Plan()

		if err != nil {
			log.Printf("[ERROR] Error reading plan: %s", err)
		}

		changes := plan.ResourceChanges()
		prj.PendingChanges = make([]model.ResourceChange, len(changes))
		for i, v := range changes {
			prj.PendingChanges[i] = *v
		}
	}

	// commit updates to the project
	err := p.store.Update(projectNS, prj.GUID, prj)
	if err != nil {
		log.Printf("[ERROR] Error updating project status: %s", err)
	}
}
