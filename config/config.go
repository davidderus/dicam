package config

import (
	"encoding/json"
	"os"
)

type camera struct {
	path      string
	role      string
	autostart bool
	notifiers []notifier
	watcher
}

type notifier struct {
	service    string
	recipients []string
}

type controller struct {
	port int
}

type watcher struct {
	autostart string
	countdown int
}

type Options struct {
	motionPath string
	controller
	cameras []camera
}

func Read(filename string) (*Options, error) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	options := &Options{}

	decodeError := decoder.Decode(options)
	if decodeError != nil {
		return nil, decodeError
	}

	return options, nil
}
