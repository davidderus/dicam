package server

import (
	"net/http"
	"path/filepath"
)

// HomeIndex gives dicam status infos
func HomeIndex(w http.ResponseWriter, r *http.Request) {
	writeWithTemplate(w, "HomeIndex", filepath.Join("index.html"), nil)
}
