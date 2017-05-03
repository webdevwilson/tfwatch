package routes

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"os"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/terraform"
	"github.com/webdevwilson/terraform-ci/config"
	"github.com/webdevwilson/terraform-ci/model"
)

func init() {
	r := Router()
	r.HandleFunc("/api/projects/{guid}/tfplan", wrapHandler(planGet)).Methods("GET")
	r.HandleFunc("/api/projects/{guid}/tfplan", wrapHandler(planApply)).Methods("POST")
}

type planDescription struct {
	Resources []resourceChange `json:"resources"`
}

type resourceChange struct {
	Name   string `json:"name"`
	Action string `json:"action"`
}

func planGet(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]
	dir := config.Get().CheckoutDirectory

	project, err := model.GetProject(guid)
	if err != nil {
		return
	}

	plan, err := os.Open(path.Join(dir, project.Name, "terraform.tfplan"))
	if err != nil {
		return
	}

	tfPlan, err := terraform.ReadPlan(plan)
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

	return
}

func planApply(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]
	data, err = model.ExecutePlan(guid)
	return
}

func resourceCount(plan *terraform.Plan) int {
	return 0
}
