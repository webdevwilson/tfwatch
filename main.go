package main

import (
	"github.com/webdevwilson/terraform-ci/config"
	_ "github.com/webdevwilson/terraform-ci/execute"
)

func main() {
	settings := config.Get()

	go settings.Server.Start()

	// loop
	for {
	}
}
