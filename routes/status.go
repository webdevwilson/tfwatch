package routes

import "net/http"

func init() {
	registrationCh <- func(s *server) {
		s.registerAPIEndpoints(api{"GET", "/status", status})
	}
}

// Status returns error
func status(*http.Request) (interface{}, error) {
	return "OK", nil
}
