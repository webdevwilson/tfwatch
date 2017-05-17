package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"time"

	"github.com/hashicorp/logutils"
	"github.com/webdevwilson/terraform-ci/controller"
	"github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/persist"
	"github.com/webdevwilson/terraform-ci/routes"
)

// Settings contains all the configuration values for the service. Do not put public stuff
// in this data structure. Should configure application here and expose interfaces, not config values (IOC).
// This will likely be discarded in favor of a 'Context' data structure that contains references to the disparate
// systems. Even further, these systems should be communicating across a messaging channel as opposed to being
// tightly coupled.
type Context struct {
	Server   routes.HTTPServer
	Projects controller.Projects
}

// Options for configuring the application
type Options struct {
	CheckoutDir string
	StateDir    string
	ClearState  bool
	LogDir      string
	LogLevel    logutils.LogLevel
	SiteDir     string
	Port        uint
	RunPlan     bool
}

// NewContext creates the execution context for server. The context is the root
// data structure containing other data structures. The context should not be called
// other than during initialization of the process.
func NewContext(opts *Options) *Context {

	// configure checkout directory and ensure it exists
	if _, err := os.Stat(opts.CheckoutDir); os.IsNotExist(err) {
		log.Fatalf("[FATAL] Checkout directory \"%s\" does not exist", opts.CheckoutDir)
	}

	// clear state directory
	if opts.ClearState {
		log.Printf("[WARN] Clearing state directory '%s'", opts.StateDir)
		err := os.RemoveAll(opts.StateDir)
		if err != nil {
			log.Printf("[WARN] Error clearing state directory '%s': %s", opts.StateDir, err)
		}
	}

	// create state directory
	err := os.MkdirAll(opts.StateDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating state directory '%s': %s", opts.StateDir, err)
	}

	// initialize the data store
	store, err := persist.NewBoltStore(path.Join(opts.StateDir))
	if err != nil {
		log.Fatalf("[FATAL] Error initializing persistence: %s", err)
	}

	// logging configuration
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel(opts.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.Printf("[INFO] Log level set to %s", opts.LogLevel)

	// create an executor
	executor := execute.NewExecutor(store, path.Join(opts.LogDir, "executor"))

	// create the controller
	projects := controller.NewProjectsController(opts.CheckoutDir, store, executor, 5*time.Minute, opts.RunPlan)

	// create the HTTP server
	accessLogDir := path.Join(opts.LogDir, "http")
	err = os.MkdirAll(accessLogDir, os.ModePerm)
	if err != nil {
		log.Printf("[WARN] Error creating access log directory '%s': %s", accessLogDir, err)
	}

	// open access log
	accessLog, err := os.OpenFile(path.Join(accessLogDir, "access.log"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[WARN] Error opening access log: %s", err)
	}

	siteDir := os.Getenv("SITE_DIR")
	port := envOrUint("PORT", 3000)
	server := routes.InitializeServer(port, accessLog, projects, siteDir)

	// initialize the context
	return &Context{
		Projects: projects,
		Server:   server,
	}
}

// env returns environment variables. fatal error if it does not exist
func env(name string) (v string) {
	if v = os.Getenv(name); len(v) == 0 {
		log.Fatalf("[FATAL] %s variable required.", name)
	}
	return
}

// envOrInt returns the int representation of the environment variable or the default
func envOrUint(name string, defaultVal uint) uint {
	var v string

	// does an environment variable exist?
	if v = os.Getenv(name); len(v) == 0 {
		return defaultVal
	}

	// convert environment variable
	val, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("[WARN] Invalid int given for PORT environment variable: %s", err)
		return defaultVal
	}

	return uint(val)
}
