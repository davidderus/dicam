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
