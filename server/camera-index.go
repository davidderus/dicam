package server

import (
	"net/http"
	"path/filepath"
	"strings"
)

// CameraIndex lists all cameras
func CameraIndex(w http.ResponseWriter, r *http.Request) {
	camsRawList := askClient("CAM-LIST")
	camsList := strings.Split(camsRawList, "\n")

	writeWithTemplate(w, "CameraIndex", filepath.Join("cameras", "index.html"), camsList)
}
