package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// CameraStart sends a request to start a given camera
func CameraStart(w http.ResponseWriter, r *http.Request) {
	// TODO prevent code duplication
	vars := mux.Vars(r)
	cameraID := vars["cameraId"]

	_, camStartError := askClient("CAM-START-" + cameraID)

	if camStartError != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, camStartError)
	} else {
		// Waiting for motion to make the camera available
		time.Sleep(750 * time.Millisecond)

		http.Redirect(w, r, "/cameras/"+cameraID, 302)
	}
}

// CameraStop sends a request to stop a given camera
func CameraStop(w http.ResponseWriter, r *http.Request) {
	// TODO prevent code duplication
	vars := mux.Vars(r)
	cameraID := vars["cameraId"]

	_, camStopError := askClient("CAM-STOP-" + cameraID)

	if camStopError != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, camStopError)
	} else {
		http.Redirect(w, r, "/cameras/"+cameraID, 302)
	}
}
