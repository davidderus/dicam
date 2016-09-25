package main

import "log"

type controller struct {
	cameras []camera
}

func (c *controller) launchCamera(cameraID int) camera {
	cam := camera{id: cameraID}

	cam.setup()

	log.Printf("Starting cam %d\n", cam.id)
	cam.start()

	log.Printf("Camera started with PID %d\n", cam.pid)

	c.cameras = append(c.cameras, cam)

	return cam
}

func (c controller) listCameras() []camera {
	return c.cameras
}

func (c controller) getCameraByID(cameraID int) *camera {
	for _, cam := range c.cameras {
		if cam.id == cameraID {
			return &cam
		}
	}

	panic("No camera found")
}

func (c controller) stopCamera(cameraID int) {
	log.Printf("Stopping cam %d\n", cameraID)

	cam := c.getCameraByID(cameraID)
	pid := cam.pid

	cam.stop()

	log.Printf("Camera stopped via PID %d\n", pid)
}

func main() {
	mainController := controller{}

	mainController.launchCamera(1)

	mainController.stopCamera(1)
}
