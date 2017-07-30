package server

import "net/http"

// HomeIndex gives dicam status infos
func HomeIndex(w http.ResponseWriter, r *http.Request) {
	// askClient("STATS")
	// writeWithTemplate(w, "HomeIndex", filepath.Join("index.html"), nil)

	// Index page is not ready yet, so we're redirecting to the cams page
	http.Redirect(w, r, "/cameras", 302)
}
