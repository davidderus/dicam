package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/davidderus/dicam/client"
	"github.com/davidderus/dicam/config"
	"github.com/gorilla/mux"
)

var AppConfig *config.Config

// Start starts the webserver on :8000
func Start() {
	router := mux.NewRouter()

	AppConfig = loadConfig()

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

// askClient interacts once with the command center
func askClient(command string) string {
	config := AppConfig
	client := &client.Client{Host: config.Host, Port: config.Port}
	connectionError := client.Connect()

	if connectionError != nil {
		log.Fatalln(connectionError)
	}

	askResponse, _ := client.Ask(command)
	return askResponse
}

// loadConfig reads the dicam config file
func loadConfig() *config.Config {
	config, readError := config.Read()
	if readError != nil {
		log.Fatalln(readError)
	}

	return config
}
