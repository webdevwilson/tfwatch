package options

import (
	"fmt"
	"log"
	"os"
	"flags"

	"github.com/hashicorp/logutils"
	"github.com/webdevwilson/tfwatch/controller"
	"github.com/webdevwilson/tfwatch/execute"
)

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

func ParseArgs(args []string) Configuration {

	var checkoutDir, logDir, logLevel, siteDir, stateDir string
	var port uint
	var clearState, help, noPlanRuns, verbose bool


	flags := flag.NewFlagSet("tfwatch", flag.ExitOnError)
	flags.BoolVar(&clearState, "clear-state", false, "Remove all state before starting")
	flags.BoolVar(&help, "h", false, "")
	flags.BoolVar(&help, "help", false, "Display usage information")
	flags.StringVar(&logDir, "log-dir", "", "Directory the logs will be placed in")
	flags.StringVar(&logLevel, "log-level", envOr("LOG_LEVEL", "INFO"), "Log level. One of DEBUG, INFO, WARN, ERROR")
	flags.BoolVar(&noPlanRuns, "no-plans", false, "Prevents tfwatch from updating the plans")
	flags.UintVar(&port, "port", 3000, "Defines port HTTP server will bind to")
	flags.StringVar(&siteDir, "site-dir", envOr("SITE_DIR", "site"), "Directory site is served from")
	flags.StringVar(&stateDir, "state-dir", envOr("STATE_DIR", ""), "Directory where state is stored")
	flags.BoolVar(&verbose, "v", false, "")
	flags.BoolVar(&verbose, "verbose", false, "Configure max logging")

	//flag.Usage = usage
	flags.Parse(os.Args[1:])

	// print helpful usage information
	if help {
		flags.Usage()
		os.Exit(0)
	}

	if verbose {
		logLevel = "DEBUG"
	}

	// ensure we have a checkout directory, this is the only required option
	if checkoutDir = flags.Arg(0); checkoutDir == "" {
		log.Printf("[ERROR] No directory specified!")
		flags.Usage()
		os.Exit(1)
	}

	// set defaults that use checkout directory
	if stateDir == "" {
		stateDir = path.Join(checkoutDir, ".tfwatch")
	}

	if logDir == "" {
		logDir = path.Join(stateDir, "logs")
	}

	logLevel = strings.ToUpper(logLevel)

	return &options.Configuration{
		CheckoutDir: checkoutDir,
		ClearState:  clearState,
		LogDir:      logDir,
		LogLevel:    logutils.LogLevel(logLevel),
		Port:        uint16(port),
		RunPlan:     !noPlanRuns,
		SiteDir:     siteDir,
		StateDir:    stateDir,
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
