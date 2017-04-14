package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/webdevwilson/terraform-ui/core/plan"
)

func init() {
	register("/api/plans", planList)
	register("/api/plans/{guid}", planGet)
	register("/api/plans/{guid}/apply", planApply)
}

// planList writes a list of plans to the request
func planList(req *http.Request) (data interface{}, err error) {
	data, err = plan.List()
	return
}

// planGet fetches a plan
func planGet(req *http.Request) (data interface{}, err error) {
	name := mux.Vars(req)["guid"]
	data, err = plan.Get(name)
	return
}

// planApply runs the designated plan
func planApply(req *http.Request) (data interface{}, err error) {
	name := mux.Vars(req)["guid"]
	data, err = plan.Apply(name)
	return
}
