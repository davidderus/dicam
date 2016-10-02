package controller

import (
	"log"

	"github.com/davidderus/dicam/config"
)

var CamsPoolInstance *CamsPool

func Start(config *config.Config) error {
	CamsPoolInstance = &CamsPool{config: config}

	camsToStart := config.ListCamsToStart()

	// Starting Cameras with Autostart to true
	for _, cameraID := range camsToStart {
		log.Printf("Autostarting camera %s", cameraID)
		output, autostartError := CamsPoolInstance.launchCamera(cameraID)
		if autostartError != nil {
			log.Println(autostartError)
		} else {
			log.Println(output)
		}
	}

	cc := &CommandCenter{Port: config.Port}
	startError := cc.Start()

	if startError != nil {
		return startError
	}

	return nil
}
