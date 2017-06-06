package test

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

// SuppressLogs should be called by tests to turn down the log level
func SuppressLogs() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}
