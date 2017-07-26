package context

import (
	"log"
	"os"
	"path"
	"time"

	"fmt"
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
type Instance struct {
	Server   routes.HTTPServer
	Projects controller.Projects
	System   controller.System
}

// Configuration settings for the application
type Configuration struct {
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
func NewContext(cfg *Configuration) *Instance {

	// configure checkout directory and ensure it exists
	if _, err := os.Stat(cfg.CheckoutDir); os.IsNotExist(err) {
		log.Fatalf("[FATAL] Checkout directory \"%s\" does not exist", cfg.CheckoutDir)
	}

	// clear state directory
	if cfg.ClearState {
		log.Printf("[WARN] Clearing state directory '%s'", cfg.StateDir)
		err := os.RemoveAll(cfg.StateDir)
		if err != nil {
			log.Printf("[WARN] Error clearing state directory '%s': %s", cfg.StateDir, err)
		}
	}

	// create state directory
	err := os.MkdirAll(cfg.StateDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating state directory '%s': %s", cfg.StateDir, err)
	}

	// initialize the data store
	store, err := persist.NewBoltStore(path.Join(cfg.StateDir))
	if err != nil {
		log.Fatalf("[FATAL] Error initializing persistence: %s", err)
	}

	// logging configuration
	configureLogging(cfg.LogLevel)

	// create an executor
	executor := execute.NewExecutor(store, path.Join(cfg.LogDir, "executor"))

	// create the controller
	projects := controller.NewProjectsController(cfg.CheckoutDir, store, executor, 5*time.Minute, cfg.RunPlan)

	// create the system controller
	system := controller.NewSystemController(systemConfigValues(cfg), executor)

	// create the HTTP server
	accessLogDir := path.Join(cfg.LogDir, "http")
	err = os.MkdirAll(accessLogDir, os.ModePerm)
	if err != nil {
		log.Printf("[WARN] Error creating access log directory '%s': %s", accessLogDir, err)
	}

	// open access log
	accessLog, err := os.OpenFile(path.Join(accessLogDir, "access.log"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[WARN] Error opening access log: %s", err)
	}

	siteDir := cfg.SiteDir
	port := cfg.Port
	server := routes.InitializeServer(port, accessLog, system, projects, siteDir)

	// initialize the context
	return &Instance{
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

func systemConfigValues(cfg *Configuration) []controller.SystemConfigurationValue {
	return []controller.SystemConfigurationValue{
		{"CheckoutDir", "Checkout Directory", cfg.CheckoutDir},
		{"LogLevel", "Log Level", string(cfg.LogLevel)},
		{"Port", "HTTP Port", fmt.Sprintf("%d", cfg.Port)},
	}
}
