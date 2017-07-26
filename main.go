package main

import (
	"os"

	"github.com/webdevwilson/tfwatch/context"
	_ "github.com/webdevwilson/tfwatch/execute"
	"github.com/webdevwilson/tfwatch/options"
)

func main() {

	cfg := options.ParseArgs(os.Args[1:])
	ctx := context.NewContext(cfg)

	go ctx.Server.Start()

	// loop
	for {
	}
}

// envOr returns the environment variable or the default values
func envOr(name string, defaultVal string) (v string) {
	if v = os.Getenv(name); len(v) == 0 {
		v = defaultVal
	}
	return
}
