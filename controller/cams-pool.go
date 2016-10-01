package controller

import (
	"fmt"
	"log"
	"strings"
)

type CamsPool struct {
	cameras []*camera
}

func (cp *CamsPool) launchCamera(cameraID int) *camera {
	cam := &camera{id: cameraID}

	setupError := cam.setup()

	if setupError != nil {
		log.Fatalln("Error during camera setup:", setupError)
	} else {
		log.Printf("Starting cam %d\n", cam.id)
	}

	startError := cam.start()

	if startError != nil {
		log.Fatalln("Error during camera launch:", startError)
	} else {
		log.Printf("Camera started with PID %d\n", cam.pid)
	}

	cp.cameras = append(cp.cameras, cam)

	return cam
}

func (cp *CamsPool) listCameras() string {
	cams := cp.cameras
	message := "No camera"

	if len(cams) > 0 {
		var camsList []string

		for _, cam := range cams {
			camsList = append(camsList, fmt.Sprintf("Cam. %d - PID %d", cam.id, cam.pid))
		}

		message = strings.Join(camsList, "\n")
	}

	return message
}

func (cp *CamsPool) getCameraByID(cameraID int) *camera {
	for _, cam := range cp.cameras {
		if cam.id == cameraID {
			return cam
		}
	}

	log.Fatalln("No camera found")

	return nil
}

func (cp *CamsPool) stopCamera(cameraID int) {
	log.Printf("Stopping cam %d\n", cameraID)

	cam := cp.getCameraByID(cameraID)
	pid := cam.pid

	err := cam.stop()

	if err != nil {
		log.Fatalln("Error while stopping camera:", err)
	} else {
		log.Printf("Camera stopped via PID %d\n", pid)
	}
}
