package controller

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/davidderus/dicam/config"
)

const streamPortStart = 8780

// CamsPool stores all started Cams. It also allows easy access to the global
// config
type CamsPool struct {
	cameras           []*camera
	config            *config.Config
	streamPortMapping map[string]int
	currentFreePort   int
}

// launchCamera assures a given camera setup and launch
func (cp *CamsPool) launchCamera(cameraID string) (string, error) {
	cam := &camera{ID: cameraID}
	camOptions, cameraOptionsError := cp.config.GetCameraOptions(cameraID)
	if cameraOptionsError != nil {
		return "", cameraOptionsError
	}

	cam.setWorkingDir(cp.config.WorkingDir)

	streamPort := cp.allocateStreamPort(cameraID)
	cam.setStreamPort(streamPort)

	setupError := cam.setup(camOptions)

	if setupError != nil {
		return "", setupError
	}

	startError := cam.start()

	if startError != nil {
		return "", startError
	}

	cp.cameras = append(cp.cameras, cam)

	return fmt.Sprintf("Camera %s started with PID %d\n", cam.ID, cam.pid), nil
}

// listCameras return all the config cameras.
// If the camera is running, its PID is also returned
func (cp *CamsPool) listCameras() (string, error) {
	cams := cp.cameras
	message := "No camera"
	camsList := []string{}
	camIDS := []string{}

	// Listing running cams first
	for _, runningCam := range cams {
		camsList = append(camsList, fmt.Sprintf("Cam. %s - PID %d", runningCam.ID, runningCam.pid))
		camIDS = append(camIDS, runningCam.ID)
	}

	for camName := range cp.config.Cameras {
		if inSlice(camName, camIDS) {
			continue
		}
		camsList = append(camsList, fmt.Sprintf("Cam. %s - Not running", camName))
	}

	if len(camsList) > 0 {
		message = strings.Join(camsList, "\n")
	}

	return message, nil
}

// inSlice indicates if a string is available in an array of strings
func inSlice(needle string, haystack []string) bool {
	index := sort.SearchStrings(haystack, needle)
	return index < len(haystack)
}

// getCameraByID returns a Camera instance from the CamsPool
func (cp *CamsPool) getCameraByID(cameraID string) (*camera, error) {
	for _, cam := range cp.cameras {
		if cam.ID == cameraID {
			return cam, nil
		}
	}

	return nil, fmt.Errorf("No camera %s found running", cameraID)
}

// stopCamera stops a camera from the CamsPool
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

func (cp *CamsPool) allocateStreamPort(cameraID string) int {
	streamPort := cp.currentFreePort

	cp.currentFreePort = streamPort + 1

	cp.streamPortMapping[cameraID] = streamPort

	return streamPort
}

// boot initiates the CamsPool by launching all autostarting Cameras
//
// TODO Improve message code logging
// TODO Externalize logging too
func (cp *CamsPool) boot() {
	log.Printf("Cameras Pool working dir: %s", cp.config.WorkingDir)

	// Initiates currentFreePort
	cp.currentFreePort = streamPortStart
	cp.streamPortMapping = make(map[string]int)

	// Lists cams to start (auto_start to true)
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

// getCameraInfos returns useful informations about a camera
func (cp *CamsPool) getCameraInfos(cameraID string) (string, error) {
	cam, findError := cp.getCameraByID(cameraID)
	if findError != nil {
		return "", findError
	}

	infos := cam.infos()

	return fmt.Sprintf(infos), nil
}
