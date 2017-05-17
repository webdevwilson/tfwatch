//go:generate esc -o site_static.go -private -ignore .map -pkg routes ../site/dist

package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	registrationCh <- func(s *server) {
		// _escFS is generated in site_static.go from source in site/dist
		fs := http.FileServer(http.Dir(s.siteDir))

		s.registerEndpoint("GET", "/", redirectRoot)
		s.registerEndpoint("GET", "/site/{path:.*}", prefix("/site", fs))
		s.registerEndpoint("GET", "/dist/{path:.*}", prefix("/site/dist", fs))
	}
}

// prefix adds a prefix to every request (see http.StripPrefix)
func prefix(prefix string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)["path"]
		r.URL.Path = fmt.Sprintf("%s/%s", prefix, path)
		log.Printf("[INFO] routes.prefix() Path: '%s', rewritten to '%s'", path, r.URL.Path)
		h.ServeHTTP(w, r)
		log.Printf("[DEBUG] routes.prefix() h.ServeHTTP() complete")
	}
}

// redirectRoot sends a redirect for root url requests
func redirectRoot(resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, "/site/index.html#", http.StatusTemporaryRedirect)
}
