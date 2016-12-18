package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Start starts the webserver on :8000
func Start() {
	router := mux.NewRouter()

	router.HandleFunc("/", HomeIndex)
	router.HandleFunc("/cameras", CameraIndex)
	router.HandleFunc("/cameras/{cameraId}", CameraShow)

	http.Handle("/", router)

	server := &http.Server{
		Handler: router,
		Addr:    ":8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

// HomeIndex gives dicam status infos
func HomeIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Index")
}

// CameraIndex lists all cameras
func CameraIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Cameras index")
}

// CameraShow shows a given camera
func CameraShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["cameraId"]

	fmt.Fprintln(w, "Showing camera "+cameraID)
}
