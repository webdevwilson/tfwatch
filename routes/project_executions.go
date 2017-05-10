package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	registrationCh <- func(s *server) {
		s.registerAPIEndpoints([]api{
			api{"GET", "/api/projects/{guid}/executions", projectExecutions},
		}...)
	}
}

func projectExecutions(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]
	prj, err := projectsController().Get(guid)
	if err != nil {
		return
	}
	data, err = projectsController().GetExecutions(prj)
	return
}
