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
	// FS() is created by esc and returns a http.Filesystem.
	fs := _escFS(false)
	r.HandleFunc("/", redirectRoot)
	r.Handle("/site/{path:.*}", prefix("/site/dist", http.FileServer(fs))).Methods("GET")
	r.Handle("/static/{path:.*}", prefix("/site/dist/static", http.FileServer(fs))).Methods("GET")
}

// // prefix adds a prefix to every request see http.StripPrefix
func prefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)["path"]
		log.Printf("[INFO] %s", r.URL.Path)
		r.URL.Path = fmt.Sprintf("%s/%s", prefix, path)
		log.Printf("[INFO] %s", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func redirectRoot(resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, "/site/projectList.html#", http.StatusPermanentRedirect)
}
