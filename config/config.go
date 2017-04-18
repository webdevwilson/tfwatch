package config

import (
	"log"
	"os"
	"path"

	"github.com/hashicorp/logutils"
	"github.com/webdevwilson/terraform-ui/persist"
	"github.com/webdevwilson/terraform-ui/task"
)

// Settings contains all the configuration values for the service
type Settings struct {
	LogLevel string
	SiteRoot string
	Port     string
	Store    persist.Store
	Executor *task.Executor
}

var settings *Settings

func init() {

	// create the persistent store
	store, err := persist.NewLocalFileStore(envOrFunc("STATE_PATH", defaultStatePath))

	if err != nil {
		log.Fatalf("[FATAL] Error initializing persistence: %s", err)
	}

	executor := task.NewExecutor(store)

	// initialize settings
	settings = &Settings{
		envOr("LOG_LEVEL", "INFO"),
		envOrFunc("SITE_ROOT", defaultSiteRoot),
		envOr("PORT", "3000"),
		store,
		executor,
	}

	log.Printf("[INFO] Log level set to %s", settings.LogLevel)

	// configure logging
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel(settings.LogLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
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
	statePath = path.Join(wd, ".state")
	return
}

func defaultSiteRoot() (siteRoot string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[FATAL] Failed to get current working directory: %s", wd)
	}
	siteRoot = path.Join(wd, "site", "dist")
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

// envOrFunc returns the named environment value or the result of executing the function
func envOrFunc(name string, defaultFunc func() string) string {
	var v string
	if v = os.Getenv(name); len(v) == 0 {
		v = defaultFunc()
	}
	return v
}
