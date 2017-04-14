package routes

import (
	"net/http"

	"github.com/webdevwilson/terraform-ui/core/project"
)

func init() {
	register("/api/projects", projectList)
}

func projectList(req *http.Request) (data interface{}, err error) {
	data, err = project.List()
	return
}
