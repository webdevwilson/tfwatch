//go:generate esc -o site_static.go -private -ignore .map -pkg routes ../site/dist

package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := Router()
	// _escFS is generated in site_static.go from source in site/dist
	fs := http.FileServer(newFS())
	r.HandleFunc("/", redirectRoot)
	r.Handle("/site/{path:.*}", prefix("/site/dist", fs)).Methods("GET")
	r.Handle("/static/{path:.*}", prefix("/site/dist/static", fs)).Methods("GET")
}

// prefix adds a prefix to every request (see http.StripPrefix)
func prefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)["path"]
		r.URL.Path = fmt.Sprintf("%s/%s", prefix, path)
		h.ServeHTTP(w, r)
	})
}

// redirectRoot sends a redirect for root url requests
func redirectRoot(resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, "/site/index.html#", http.StatusTemporaryRedirect)
}

// A simple wrapper that logs errors
type fsWrapper struct {
	handler http.FileSystem
}

// newFS creates an FS that wraps the static generated http.FileSystem and logs when errors occur
func newFS() http.FileSystem {
	return fsWrapper{_escFS(false)}
}

func (fs fsWrapper) Open(name string) (file http.File, err error) {
	file, err = fs.handler.Open(name)

	if err != nil {
		log.Printf("[ERROR] Error serving file '%s': %s", name, err)
	}

	return
}
