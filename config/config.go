package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Camera struct {
	ID        string
	Path      string
	Role      string
	Autostart bool
	Notifiers []*Notifier
	*Watcher
}

type Notifier struct {
	Service    string
	Recipients []string
}

type Controller struct {
	Port int
}

type Watcher struct {
	Autostart string
	Countdown int
}

type Options struct {
	MotionPath string
	*Controller
	Cameras []*Camera
}

func Read(filename string) (*Options, error) {
	file, _ := os.Open(filename)
	defer file.Close()

	decoder := json.NewDecoder(file)
	options := &Options{}

	decodeError := decoder.Decode(options)
	if decodeError != nil {
		return nil, decodeError
	}

	return options, nil
}

func (o *Options) GetAutostartCameras() []*Camera {
	autostartCameras := []*Camera{}

	for _, cam := range o.Cameras {
		if cam.Autostart == true {
			autostartCameras = append(autostartCameras, cam)
		}
	}

	return autostartCameras
}

func (o *Options) GetCameraByID(cameraID string) (*Camera, error) {
	for _, cam := range o.Cameras {
		if cam.ID == cameraID {
			return cam, nil
		}
	}

	return nil, fmt.Errorf("No camera %s found", cameraID)
}
