package routes

import (
	"net/http"
	"path"

	"os"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/terraform"
	"github.com/webdevwilson/terraform-ci/config"
	"github.com/webdevwilson/terraform-ci/model"
)

func init() {
	r := Router()
	r.HandleFunc("/api/projects/{guid}/tfplan", wrapHandler(planGet)).Methods("GET")
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

	data, err = terraform.ReadPlan(plan)

	return
}
