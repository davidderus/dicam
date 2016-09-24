package main

import (
	"fmt"
	"time"
)

type controller struct {
	cameras []camera
}

func (c *controller) launchCamera(cameraID int) camera {
	cam := camera{id: cameraID}

	c.cameras = append(c.cameras, cam)

	cam.setup()

	fmt.Printf("Starting cam %d\n", cam.id)
	cam.start()

	fmt.Printf("Camera started with PID %d\n", cam.command.Process.Pid)

	return cam
}

func main() {
	mainController := controller{}

	for i := 1; i < 6; i++ {
		fmt.Println("Launching camera", i)
		mainController.launchCamera(i)
	}

	time.Sleep(10)

	for i := 1; i < 4; i++ {
		fmt.Println("Stopping camera", i)
		mainController.cameras[i].stop()
	}

	time.Sleep(10)

	for i := 4; i < 6; i++ {
		fmt.Println("Stopping camera", i)
		mainController.cameras[i].stop()
	}
}
