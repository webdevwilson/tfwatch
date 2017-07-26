package main

import (
	"flag"
	"github.com/hashicorp/logutils"
	"github.com/webdevwilson/tfwatch/context"
	"log"
	"os"
	"path"
	"strings"
)

func main() {

	cfg := ParseArgs(os.Args[1:])
	ctx := context.NewContext(cfg)

	go ctx.Server.Start()

	// loop
	for {
	}
}

func ParseArgs(args []string) *context.Configuration {

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

	return &context.Configuration{
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

// envOr returns the environment variable or the default values
func envOr(name string, defaultVal string) (v string) {
	if v = os.Getenv(name); len(v) == 0 {
		v = defaultVal
	}
	return
}
