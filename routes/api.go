package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

// function contract for API endpoints used to add common behavior to all API endpoints
type apiHandlerFunc func(*http.Request) (interface{}, error)

type api struct {
	method  string
	path    string
	handler apiHandlerFunc
}

// ServeHTTP is used to service all API requests
func (api apiHandlerFunc) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] %s %s", req.Method, req.URL)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Handler recovered from panic: %s\n%s", r, debug.Stack())
		}
	}()

	data, err := api(req)
	if err != nil {
		log.Printf("[ERROR] Error in handler '%s': %s", req.URL, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
	resp.Header().Add("Access-Control-Allow-Credentials", "true")
	err = json.NewEncoder(resp).Encode(data)
	if err != nil {
		log.Printf("[ERROR] Error encoding response '%s': %s", req.URL, err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}
