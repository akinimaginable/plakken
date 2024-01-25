package static

import (
	"net/http"
)

// ServeStatic Serve static file from static
func ServeStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/"+r.PathValue("file")) // TODO: vérifier si c'est safe
}

// Home Serve index.html
func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}
