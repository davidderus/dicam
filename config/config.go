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

type options struct {
	motionPath string
	controller
	cameras []camera
}

type Config struct {
	*options
}

func (c *Config) read(filename string) (*Config, error) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	config := &Config{}
	options := &options{}

	decodeError := decoder.Decode(options)
	if decodeError != nil {
		return nil, decodeError
	}

	config.options = options

	return config, nil
}
