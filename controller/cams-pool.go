package controller

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/davidderus/dicam/config"
)

type CamsPool struct {
	cameras []*camera
	config  *config.Config
}

func (cp *CamsPool) launchCamera(cameraID string) (string, error) {
	cam := &camera{id: cameraID}
	camOptions, cameraOptionsError := cp.config.GetCameraOptions(cameraID)
	if cameraOptionsError != nil {
		return "", cameraOptionsError
	}

	setupError := cam.setup(camOptions)

	if setupError != nil {
		return "", setupError
	}

	startError := cam.start()

	if startError != nil {
		return "", startError
	}

	cp.cameras = append(cp.cameras, cam)

	return fmt.Sprintf("Camera %s started with PID %d\n", cam.id, cam.pid), nil
}

func (cp *CamsPool) listCameras() (string, error) {
	cams := cp.cameras
	message := "No camera"
	camsList := []string{}
	camIDS := []string{}

	// Listing running cams first
	for _, runningCam := range cams {
		camsList = append(camsList, fmt.Sprintf("Cam. %s - PID %d", runningCam.id, runningCam.pid))
		camIDS = append(camIDS, runningCam.id)
	}

	for camName := range cp.config.Cameras {
		if inSlice(camName, camIDS) {
			continue
		}
		camsList = append(camsList, fmt.Sprintf("Cam. %s - Not running", camName))
	}

	message = strings.Join(camsList, "\n")

	return message, nil
}

func inSlice(needle string, haystack []string) bool {
	index := sort.SearchStrings(haystack, needle)
	return index < len(haystack)
}

func (cp *CamsPool) getCameraByID(cameraID string) (*camera, error) {
	for _, cam := range cp.cameras {
		if cam.id == cameraID {
			return cam, nil
		}
	}

	return nil, fmt.Errorf("No camera %s found", cameraID)
}

func (cp *CamsPool) stopCamera(cameraID string) (string, error) {
	cam, findError := cp.getCameraByID(cameraID)
	if findError != nil {
		return "", findError
	}

	err := cam.stop()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Camera %s stopped via PID %d\n", cameraID, cam.pid), nil
}

// @todo Improve message code logging
// @todo Externalize logging too
func (cp *CamsPool) boot() {
	camsToStart := cp.config.ListCamsToStart()

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
}
