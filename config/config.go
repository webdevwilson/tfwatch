package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"time"

	"github.com/hashicorp/logutils"
	"github.com/webdevwilson/tfwatch/controller"
	"github.com/webdevwilson/tfwatch/execute"
	"github.com/webdevwilson/tfwatch/persist"
	"github.com/webdevwilson/tfwatch/routes"
)

// Settings contains all the configuration values for the service. Do not put public stuff
// in this data structure. Should configure application here and expose interfaces, not config values (IOC).
// This will likely be discarded in favor of a 'Context' data structure that contains references to the disparate
// systems. Even further, these systems should be communicating across a messaging channel as opposed to being
// tightly coupled.
type Context struct {
	Server   routes.HTTPServer
	Projects controller.Projects
	System   controller.System
}

// Options for configuring the application
type Options struct {
	CheckoutDir string
	StateDir    string
	ClearState  bool
	LogDir      string
	LogLevel    logutils.LogLevel
	SiteDir     string
	Port        uint16
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
	configureLogging(opts.LogLevel)

	// create an executor
	executor := execute.NewExecutor(store, path.Join(opts.LogDir, "executor"))

	// create the controller
	projects := controller.NewProjectsController(opts.CheckoutDir, store, executor, 5*time.Minute, opts.RunPlan)

	// create the system controller
	system := systemController(opts, executor)

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

	siteDir := opts.SiteDir
	port := opts.Port
	server := routes.InitializeServer(port, accessLog, system, projects, siteDir)

	// initialize the context
	return &Context{
		Projects: projects,
		Server:   server,
		System:   system,
	}
}

func configureLogging(level logutils.LogLevel) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: level,
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.Printf("[INFO] Log level set to %s", level)
}

func systemController(opts *Options, executor execute.Executor) controller.System {
	config := make([]*controller.SystemConfigurationValue, 3)
	config[0] = cfg("CheckoutDir", "Checkout Directory", opts.CheckoutDir)
	config[1] = cfg("LogLevel", "Log Level", string(opts.LogLevel))
	config[2] = cfg("Port", "HTTP Port", fmt.Sprintf("%d", opts.Port))
	return controller.NewSystemController(config, executor)
}

func cfg(id, name, value string) *controller.SystemConfigurationValue {
	return &controller.SystemConfigurationValue{
		ID:    id,
		Name:  name,
		Value: value,
	}
}
