package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/terraform"
	"github.com/webdevwilson/tfwatch/model"
)

func init() {
	registrationCh <- func(s *server) {
		s.registerAPIEndpoints([]api{
			api{"GET", "/api/projects/{guid}/tfplan", projectPlanGet},
			api{"POST", "/api/projects/{guid}/tfplan", projectPlanApply},
		}...)
	}
}

type planDescription struct {
	Resources []model.ResourceChange `json:"resources"`
}

func projectPlanGet(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]

	project, err := projectsController().Get(guid)
	if err != nil {
		return
	}

	data = planDescription{project.PendingChanges}

	log.Printf("[DEBUG] Found %d resource modifications", len(project.PendingChanges))

	return
}

func projectPlanApply(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]

	project, err := projectsController().Get(guid)
	if err != nil {
		return
	}

	data, err = projectsController().ExecutePlan(project)
	return
}

func resourceCount(plan *terraform.Plan) int {
	return 0
}
