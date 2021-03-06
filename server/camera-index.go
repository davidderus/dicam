package server

import (
	"net/http"
	"path/filepath"
)

// CameraIndex lists all cameras
func CameraIndex(w http.ResponseWriter, r *http.Request) {
	camsList := getCameras()

	writeWithTemplate(w, "CameraIndex", filepath.Join("cameras", "index.html"), camsList)
}
