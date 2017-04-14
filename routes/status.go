package routes

import "net/http"

func init() {
	register("/status", status)
}

// Status returns error
func status(*http.Request) (interface{}, error) {
	return "OK", nil
}
