package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/logutils"
	"github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/persist"
)

// Settings contains all the configuration values for the service
type Settings struct {
	LogLevel            string
	LogDir              string
	SiteRoot            string
	CheckoutDirectory   string
	Port                int
	PlanIntervalMinutes int
	Store               persist.Store
	Executor            *execute.Executor
}

var settings *Settings

func init() {

	// create the persistent store
	stateDir := envOrFunc("STATE_DIR", defaultStatePath)
	store, err := persist.NewLocalFileStore(stateDir)

	if err != nil {
		log.Fatalf("[FATAL] Error initializing persistence: %s", err)
	}

	// create an executor
	logDir := envOr("LOG_DIR", path.Join(stateDir, "logs"))
	executor := execute.NewExecutor(store, path.Join(logDir, "executor"))

	workingDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("[FATAL] Failed to get current working directory: %s", workingDir)
	}

	// initialize settings
	settings = &Settings{
		envOr("LOG_LEVEL", "INFO"),
		logDir,
		envOr("SITE_ROOT", path.Join(workingDir, "site", "dist")),
		envOr("CHECKOUT_DIR", path.Join(stateDir, "projects")),
		envOrInt("PORT", 3000),
		envOrInt("PLAN_INTERVAL", 5),
		store,
		executor,
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel(settings.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.Printf("[INFO] Log level set to %s", settings.LogLevel)
}

// Get returns configuration data
func Get() *Settings {
	return settings
}

// defaultStatePath returns the current working directory / .state
func defaultStatePath() (statePath string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[FATAL] Failed to get current working directory: %s", wd)
	}
	statePath = path.Join(wd, ".terraform-ci-data")
	return
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
func envOrInt(name string, defaultVal int) int {
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

	return val
}

// envOrFunc returns the named environment value or the result of executing the function
func envOrFunc(name string, defaultFunc func() string) string {
	var v string
	if v = os.Getenv(name); len(v) == 0 {
		v = defaultFunc()
	}
	return v
}
