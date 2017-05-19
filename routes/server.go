package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/webdevwilson/terraform-ci/controller"
)

// HTTPServer
type HTTPServer interface {
	Start()
}

type server struct {
	port      uint
	accessLog io.Writer
	projects  controller.Projects
	router    *mux.Router
	system    controller.System
	siteDir   string
}

var serverSingleton struct {
	instance *server
	init     sync.Once
}

var registrationCh = make(chan func(*server), 100)

// convenience method for getting the controller
func projectsController() controller.Projects {
	return serverSingleton.instance.projects
}

// InitializeServer creates an HTTPServer
func InitializeServer(port uint, accessLog io.Writer, system controller.System, projects controller.Projects, siteDir string) HTTPServer {
	serverSingleton.init.Do(func() {
		serverSingleton.instance = &server{
			port:      port,
			accessLog: accessLog,
			projects:  projects,
			router:    mux.NewRouter(),
			siteDir:   siteDir,
			system:    system,
		}
	})

	return serverSingleton.instance
}

// StartServer starts off the HTTP server
func (s *server) Start() {

	// register endpoints
	for register := range registrationCh {
		register(s)
		if len(registrationCh) == 0 {
			close(registrationCh)
			break
		}
	}

	// build a 'chain' of handlers wrapping the router:
	//   LoggingHandler for logging requests
	//	 RecoverHandler for recovering from panic
	handle := handlers.RecoveryHandler()(s.router)

	// if access log is enabled, add the handler
	if s.accessLog != nil {
		log.Printf("[DEBUG] Configuring access log")
		handle = handlers.CombinedLoggingHandler(s.accessLog, handle)
	}

	// Start the HTTP Server
	log.Printf("[INFO] Starting server on port %d", s.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), handle)

	if err != nil {
		log.Fatal("[FATAL] ListenAndServe error: ", err)
	}
}

// registerEndpoint binds an HTTP endpoint to the server
func (s *server) registerEndpoint(method string, path string, handler http.HandlerFunc) {
	log.Printf("[DEBUG] Registering endpoint '%s %s'", method, path)
	s.router.HandleFunc(path, handler).Methods(method)
}

// registerEndpoint binds an API endpoint to the server. API endpoints are wrapped to exhibit
// similar behavior
func (s *server) registerAPIEndpoints(endpoints ...api) {
	for _, endpoint := range endpoints {
		log.Printf("[DEBUG] Registering API endpoint '%s %s'", endpoint.method, endpoint.path)
		s.router.Handle(endpoint.path, endpoint.handler).Methods(endpoint.method)
	}
}
