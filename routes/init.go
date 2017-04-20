package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/gorilla/mux"
	"github.com/webdevwilson/terraform-ci/config"
	"github.com/webdevwilson/terraform-ci/persist"
)

// Route function interface
type Route func(*http.Request) (interface{}, error)

var routerSingleton struct {
	instance *mux.Router
	init     sync.Once
}

var storeSingleton struct {
	instance persist.Store
	init     sync.Once
}

// Router returns the router to map to
func Router() *mux.Router {
	routerSingleton.init.Do(func() {
		log.Printf("[DEBUG] Initializing router")
		routerSingleton.instance = mux.NewRouter()
	})
	return routerSingleton.instance
}

// Store returns the store
func Store() persist.Store {
	storeSingleton.init.Do(func() {
		storeSingleton.instance = config.Get().Store
	})
	return storeSingleton.instance
}

// register a handler
func register(path string, handler Route) (r *mux.Router) {
	log.Printf("[DEBUG] Registering route %s", path)
	r = Router()
	r.HandleFunc(path, wrapHandler(handler))
	return
}

// wrapHandler wraps handler functions
func wrapHandler(handler Route) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("[INFO] %s %s", req.Method, req.URL)

		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] Handler recovered from panic: %s\n%s", r, debug.Stack())
			}
		}()

		data, err := handler(req)
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
}
