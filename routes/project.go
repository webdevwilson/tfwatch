package routes

import (
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/webdevwilson/terraform-ci/model"
)

func init() {
	r := Router()
	r.HandleFunc("/api/projects", wrapHandler(projectList)).Methods("GET")
	r.HandleFunc("/api/projects/{guid}", wrapHandler(projectGet)).Methods("GET")
	r.HandleFunc("/api/projects", wrapHandler(projectCreate)).Methods("PUT")
	r.HandleFunc("/api/projects/{guid}", wrapHandler(projectUpdate)).Methods("POST")
	r.HandleFunc("/api/projects/{guid}", wrapHandler(projectDelete)).Methods("DELETE")
}

func projectList(req *http.Request) (data interface{}, err error) {
	return model.ListProjects()
}

func projectGet(req *http.Request) (interface{}, error) {
	guid := mux.Vars(req)["guid"]
	return model.GetProject(guid)
}

func projectCreate(req *http.Request) (data interface{}, err error) {
	var prj model.Project
	json.NewDecoder(req.Body).Decode(&prj)
	err = model.CreateProject(&prj)

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

	err = model.UpdateProject(&prj)

	if err != nil {
		return
	}

	data = prj

	return
}

func projectDelete(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]
	err = model.DeleteProject(guid)
	return
}
