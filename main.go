package main

import (
	"log"

	"fmt"

	"github.com/webdevwilson/terraform-ui/config"
	"github.com/webdevwilson/terraform-ui/routes"
	_ "github.com/webdevwilson/terraform-ui/task"
)

const port = 3000

func main() {

	log.Print(fmt.Sprintf("[INFO] Listening on port %d ...", port))

	routes.StartServer(config.Get().Port)

	for {
	}
}
