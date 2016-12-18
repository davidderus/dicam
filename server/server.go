package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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
	writeWithTemplate(w, "HomeIndex", filepath.Join("index.html"), nil)
}

// CameraIndex lists all cameras
func CameraIndex(w http.ResponseWriter, r *http.Request) {
	writeWithTemplate(w, "CameraIndex", filepath.Join("cameras", "index.html"), nil)
}

// CameraShow shows a given camera
func CameraShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cameraID := vars["cameraId"]

	writeWithTemplate(w, "CameraShow", filepath.Join("cameras", "show.html"), cameraID)
}

func writeWithTemplate(response http.ResponseWriter, templateName string, templatePath string, data interface{}) {
	getTemplateDirForFile := func(file string) string {
		return filepath.Join("server", "templates", file)
	}

	templateFile, parseError := template.ParseFiles(
		getTemplateDirForFile("layout.html"),
		getTemplateDirForFile("navbar.html"),
		getTemplateDirForFile(templatePath),
	)

	if parseError != nil {
		log.Fatalf("Can't parse template for %s", templateName)
	}

	templateFile.ExecuteTemplate(response, "layout", data)
}
