package controller

import (
	"log"
	"time"

	"github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/model"
)

// schedulePlan schedules a plan to run in a project, then waits and schedules it again
func (p *projects) schedulePlan(interval time.Duration, prj *model.Project) {
	time.AfterFunc(interval, func() {
		<-p.runPlan(prj)
		p.schedulePlan(interval, prj)
	})
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

	done = make(chan bool, 1)

	// when task is complete, update the project
	go func() {
		r := <-ch

		// update project
		prj.PlanUpdated = time.Now()
		switch r.ExitCode {
		case 0:
			prj.Status = "ok"
		case 2:
			prj.Status = "pending"
		default:
			prj.Status = "error"
		}
		log.Printf("[INFO] Project '%s' plan complete, updating status to '%s'", prj.GUID, prj.Status)
		err = p.store.Update(projectNS, prj.GUID, prj)
		if err != nil {
			log.Printf("[ERROR] Error updating project status: %s", err)
		}
	}()

	return
}
