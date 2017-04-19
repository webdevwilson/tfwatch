package routes

import (
	"fmt"
	"log"
	"net/http"
)

// StartServer starts off the HTTP server
func StartServer(port int) {
	log.Printf("[INFO] Starting server on port %d", port)

	// Start the HTTP Server
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), Router())
	if err != nil {
		log.Fatal("[FATAL] ListenAndServe error: ", err)
	}
}
