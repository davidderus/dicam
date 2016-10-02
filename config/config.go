package config

import (
	"encoding/json"
	"os"
)

type Camera struct {
	Path      string
	Role      string
	Autostart bool
	Notifiers []Notifier
	Watcher
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
	Controller
	Cameras []Camera
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
