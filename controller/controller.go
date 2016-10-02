package controller

import (
	"log"

	"github.com/davidderus/dicam/config"
)

var CamsPoolInstance *CamsPool

// todo Improve message code logging
func Start(config *config.Config) error {
	CamsPoolInstance = &CamsPool{config: config}

	camsToStart := config.ListCamsToStart()

	// Starting Cameras with Autostart to true
	for _, cameraID := range camsToStart {
		log.Printf("Autostarting camera %s", cameraID)
		output, autostartError := CamsPoolInstance.launchCamera(cameraID)
		if autostartError != nil {
			log.Printf("ERROR - %s", autostartError)
		} else {
			log.Printf("SUCCESS - %s", output)
		}
	}

	cc := &CommandCenter{Port: config.Port}
	startError := cc.Start()

	if startError != nil {
		return startError
	}

	return nil
}
