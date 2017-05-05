package routes

import (
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/webdevwilson/terraform-ci/model"
)

func init() {
	registrationCh <- func(s *server) {
		s.registerAPIEndpoints([]api{
			api{"GET", "/api/projects", projectList},
			api{"GET", "/api/projects/{guid}", projectGet},
			api{"PUT", "/api/projects", projectCreate},
			api{"POST", "/api/projects/{guid}", projectUpdate},
			api{"DELETE", "/api/projects/{guid}", projectDelete},
		}...)
	}
}

func projectList(req *http.Request) (data interface{}, err error) {
	return projectsController().List()
}

func projectGet(req *http.Request) (interface{}, error) {
	guid := mux.Vars(req)["guid"]
	return projectsController().Get(guid)
}

func projectCreate(req *http.Request) (data interface{}, err error) {
	var prj model.Project
	json.NewDecoder(req.Body).Decode(&prj)
	err = projectsController().Create(&prj)

	if err != nil {
		return
	}

	data = prj

	return
}

func projectUpdate(req *http.Request) (data interface{}, err error) {
	var prj model.Project
	json.NewDecoder(req.Body).Decode(&prj)

	// ensure the project has the same guid as in the url
	guid := mux.Vars(req)["guid"]
	prj.GUID = guid

	err = projectsController().Update(&prj)

	if err != nil {
		return
	}

	data = prj

	return
}

func projectDelete(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]
	err = projectsController().Delete(guid)
	return
}
