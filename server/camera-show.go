package server

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// CameraShow shows a given camera
func CameraShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["cameraId"]

	cameraRawInfos := askClient("CAM-INFOS-" + cameraID)
	cameraInfos := map[string]string{}

	for _, infosLine := range strings.Split(cameraRawInfos, "\n") {
		infos := strings.Split(infosLine, ":")

		cameraInfos[infos[0]] = infos[1]
	}

	writeWithTemplate(w, "CameraShow", filepath.Join("cameras", "show.html"), cameraInfos)
}
