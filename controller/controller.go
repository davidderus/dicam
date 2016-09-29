package controller

import "log"

type controller struct {
	cameras       []camera
	commandCenter CommandCenter
}

func (c *controller) launchCamera(cameraID int) camera {
	cam := camera{id: cameraID}

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

	log.Fatalln("No camera found")

	return nil
}

func (c controller) stopCamera(cameraID int) {
	log.Printf("Stopping cam %d\n", cameraID)

	cam := c.getCameraByID(cameraID)
	pid := cam.pid

	err := cam.stop()

	if err != nil {
		log.Fatalln("Error while stopping camera:", err)
	} else {
		log.Printf("Camera stopped via PID %d\n", pid)
	}
}

func (c *controller) startServer() {
	cc := CommandCenter{port: 8888}
	cc.start()

	c.commandCenter = cc
}

func main() {
	mainController := controller{}

	log.Println("Starting command center")
	mainController.startServer()
}
