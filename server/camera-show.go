package server

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// CameraShow shows a given camera
func CameraShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["cameraId"]

	writeWithTemplate(w, "CameraShow", filepath.Join("cameras", "show.html"), cameraID)
}
