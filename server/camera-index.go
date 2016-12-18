package server

import (
	"net/http"
	"path/filepath"
)

// CameraIndex lists all cameras
func CameraIndex(w http.ResponseWriter, r *http.Request) {
	// askClient("CAM-LIST")
	writeWithTemplate(w, "CameraIndex", filepath.Join("cameras", "index.html"), nil)
}
