package controller

import (
	"fmt"
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

	if len(cams) > 0 {
		var camsList []string

		for _, cam := range cams {
			camsList = append(camsList, fmt.Sprintf("Cam. %s - PID %d", cam.id, cam.pid))
		}

		message = strings.Join(camsList, "\n")
	}

	return message, nil
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
