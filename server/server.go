package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidderus/dicam/client"
	"github.com/davidderus/dicam/config"
	"github.com/gorilla/mux"

	auth "github.com/abbot/go-http-auth"
)

var AppConfig *config.Config

// Start starts the webserver on :8000
func Start() {
	router := mux.NewRouter()

	AppConfig = loadConfig()

	authenticator := auth.NewDigestAuthenticator("dicam.local", LookForSecret)

	router.HandleFunc("/", auth.JustCheck(authenticator, HomeIndex))
	router.HandleFunc("/cameras", auth.JustCheck(authenticator, CameraIndex))
	router.HandleFunc("/cameras/{cameraId}", auth.JustCheck(authenticator, CameraShow))
	router.HandleFunc("/cameras/{cameraId}/start", auth.JustCheck(authenticator, CameraStart))
	router.HandleFunc("/cameras/{cameraId}/stop", auth.JustCheck(authenticator, CameraStop))

	http.Handle("/", router)

	webServerAddress := fmt.Sprintf("%s:%d", AppConfig.WebServer.Host, AppConfig.WebServer.Port)

	server := &http.Server{
		Handler: router,
		Addr:    webServerAddress,

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

func LookForSecret(user, realm string) string {
	return ""
}

// askClient interacts once with the command center
func askClient(command string) (string, error) {
	config := AppConfig
	client := &client.Client{Host: config.Host, Port: config.Port}
	connectionError := client.Connect()

	if connectionError != nil {
		log.Fatalln(connectionError)
	}

	askResponse, askError := client.Ask(command)

	return askResponse, askError
}
