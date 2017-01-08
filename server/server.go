package server

import (
	"log"
	"net/http"
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

// loadConfig reads the dicam config file
func loadConfig() *config.Config {
	config, readError := config.Read()
	if readError != nil {
		log.Fatalln(readError)
	}

	return config
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
