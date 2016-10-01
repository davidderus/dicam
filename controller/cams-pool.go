package controller

import (
	"fmt"
	"log"
	"strings"
)

type CamsPool struct {
	cameras []*camera
}

func (cp *CamsPool) launchCamera(cameraID string) *camera {
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
		log.Printf("Camera %s started with PID %d\n", cam.id, cam.pid)
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

func (cp *CamsPool) getCameraByID(cameraID string) *camera {
	for _, cam := range cp.cameras {
		if cam.id == cameraID {
			return cam
		}
	}

	log.Fatalln("No camera found")

	return nil
}

func (cp *CamsPool) stopCamera(cameraID string) {
	log.Printf("Stopping cam %d\n", cameraID)

	cam := cp.getCameraByID(cameraID)
	pid := cam.pid

	err := cam.stop()

	if err != nil {
		log.Fatalln("Error while stopping camera:", err)
	} else {
		log.Printf("Camera %s stopped via PID %d\n", cameraID, pid)
	}
}
