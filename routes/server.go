package routes

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer(port string) {
	log.Printf("[INFO] Starting server on port %s", port)
	// Map routes
	http.Handle("/", Router())

	// Start the HTTP Server
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("[FATAL] ListenAndServe error: ", err)
	}
}
