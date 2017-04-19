package main

import (
	"log"

	"fmt"

	"github.com/webdevwilson/terraform-ui/config"
	"github.com/webdevwilson/terraform-ui/routes"
	_ "github.com/webdevwilson/terraform-ui/task"
)

func main() {
	port := config.Get().Port
	log.Print(fmt.Sprintf("[INFO] Listening on port %d ...", port))

	routes.StartServer(port)

	for {
	}
}
