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

// AppConfig contains the whole application configuration
var AppConfig *config.Config

// Authenticator contains the digest authenticator instance if needed
var Authenticator *auth.DigestAuth

// Start starts the webserver on :8000
func Start() {
	router := mux.NewRouter()

	AppConfig = loadConfig()

	if len(AppConfig.WebServer.User) > 0 {
		Authenticator = auth.NewDigestAuthenticator(AppConfig.WebServer.AuthRealm, lookForSecret)
	}

	router.HandleFunc("/", loadHandlerWithAuth(HomeIndex))
	router.HandleFunc("/cameras", loadHandlerWithAuth(CameraIndex))
	router.HandleFunc("/cameras/{cameraId}", loadHandlerWithAuth(CameraShow))
	router.HandleFunc("/cameras/{cameraId}/start", loadHandlerWithAuth(CameraStart))
	router.HandleFunc("/cameras/{cameraId}/stop", loadHandlerWithAuth(CameraStop))

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

// lookForSecret returns a password hash from config for a given existing user
func lookForSecret(user, realm string) string {
	for _, webUser := range AppConfig.WebServer.User {
		if webUser.Name == user {
			return webUser.Password
		}
	}

	return ""
}

// loadHandlerWithAuth check for any auth infos in config and use it for authentication.
func loadHandlerWithAuth(handler http.HandlerFunc) http.HandlerFunc {
	// If no auth infos are found, then no auth is set up.
	if Authenticator != nil {
		return auth.JustCheck(Authenticator, handler)
	}

	return handler
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
