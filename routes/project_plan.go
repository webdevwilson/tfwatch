package routes

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/terraform"
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
	Resources []resourceChange `json:"resources"`
}

type resourceChange struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

func projectPlanGet(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]

	project, err := projectsController().Get(guid)
	if err != nil {
		return
	}

	log.Printf("[DEBUG] Retrieving plan for project '%s'", project.GUID)

	tfPlan, err := project.Plan()

	if err != nil {
		log.Printf("[ERROR] Error retrieving plan for project '%s'", project.GUID)
		return
	}

	log.Printf("[DEBUG] Retrieved plan for project '%s'", project.GUID)
	resources := []resourceChange{}
	for _, module := range tfPlan.Diff.Modules {
		for id, res := range module.Resources {
			if change := res.ChangeType(); change != terraform.DiffNone {
				var changeStr string
				switch change {
				case terraform.DiffCreate:
					changeStr = "Create"
				case terraform.DiffDestroyCreate:
					changeStr = "Recreate"
				case terraform.DiffDestroy:
					changeStr = "Destroy"
				case terraform.DiffUpdate:
					changeStr = "Update"
				}
				name := fmt.Sprintf("%s.%s", strings.Join(module.Path, "."), id)
				resources = append(resources, resourceChange{
					name,
					changeStr,
				})
			}
		}
	}

	data = planDescription{resources}

	log.Printf("[DEBUG] Found %d resource modifications", len(resources))

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
