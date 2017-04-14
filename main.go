package main

import (
	"log"
	"net/http"

	"fmt"

	_ "github.com/webdevwilson/terraform-ui/core/task"
	"github.com/webdevwilson/terraform-ui/routes"
)

const port = "3000"

func main() {

	log.Print(fmt.Sprintf("[INFO] Listening on port %s ...", port))

	// Map routes
	http.Handle("/", routes.Router())

	// Start the HTTP Server
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("[FATAL] ListenAndServe error: ", err)
	}

	for {
	}
}
