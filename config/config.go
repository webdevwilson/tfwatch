package config

import (
	"log"
	"os"
	"path"
	"strconv"

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
type Settings struct {
	Server   routes.HTTPServer
	Projects controller.Projects
}

var settings *Settings

const defaultCheckoutDir = "/var/lib/terraform-ci"

func init() {

	checkoutDir := envOr("CHECKOUT_DIR", defaultCheckoutDir)
	stateDir := envOr("STATE_DIR", path.Join(checkoutDir, ".terraform-ci"))

	// clear state directory
	if clearState := os.Getenv("CLEAR_STATE"); len(clearState) > 0 {
		log.Printf("[WARN] Clearing state directory '%s'", stateDir)
		err := os.RemoveAll(stateDir)
		if err != nil {
			log.Printf("[WARN] Error clearing state directory '%s': %s", stateDir, err)
		}
	}

	// initialize the data store
	store, err := persist.NewLocalFileStore(path.Join(stateDir, "data"))
	if err != nil {
		log.Fatalf("[FATAL] Error initializing persistence: %s", err)
	}

	// logging configuration
	logDir := envOr("LOG_DIR", path.Join(stateDir, "logs"))
	logLevel := envOr("LOG_LEVEL", "INFO")
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.Printf("[INFO] Log level set to %s", logLevel)

	// create an executor
	executor := execute.NewExecutor(store, path.Join(logDir, "executor"))

	projects := controller.NewProjectsController(checkoutDir, store, executor)

	// create the HTTP server
	accessLogDir := path.Join(logDir, "http")
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

	// initialize settings
	settings = &Settings{
		Projects: projects,
		Server:   server,
	}
}

// Get returns configuration data
func Get() *Settings {
	return settings
}

// env returns environment variables. fatal error if it does not exist
func env(name string) (v string) {
	if v = os.Getenv(name); len(v) == 0 {
		log.Fatalf("[FATAL] %s variable required.", name)
	}
	return
}

// envOr returns the environment variable or the default values
func envOr(name string, defaultVal string) (v string) {
	if v = os.Getenv(name); len(v) == 0 {
		v = defaultVal
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

// envOrFunc returns the named environment value or the result of executing the function
func envOrFunc(name string, defaultFunc func() string) string {
	var v string
	if v = os.Getenv(name); len(v) == 0 {
		v = defaultFunc()
	}
	return v
}
