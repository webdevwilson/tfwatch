package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	register("/api/executions/{guid}", GetExecution)
}

// GetExecution returns the results of an execution
func GetExecution(req *http.Request) (data interface{}, err error) {
	guid := mux.Vars(req)["guid"]
	return task.GetResult(guid)
}
