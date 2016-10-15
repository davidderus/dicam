package config

import (
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
)

type CameraOptions struct {
	Device    string
	Role      string
	Autostart bool `toml:"auto_start"`
	Notifiers []*NotifierOptions
	Watcher   *WatcherOptions
}

type NotifierOptions struct {
	Service    string
	Recipients []string
}

type WatcherOptions struct {
	AutoStart bool `toml:"auto_start"`
	Countdown int
}

type Config struct {
	Port       int
	Host       string
	MotionPath string `toml:"motion_path"`
	Cameras    map[string]*CameraOptions
}

func Read() (*Config, error) {
	user, userError := user.Current()
	if userError != nil {
		return nil, userError
	}

	userHomeDir := user.HomeDir

	configFullPath := path.Join(userHomeDir, ".config/dicam/config.toml")

	var config Config

	_, configError := toml.DecodeFile(configFullPath, &config)
	if configError != nil {
		return nil, configError
	}

	config.setDefaults()

	validationError := config.validate()
	if validationError != nil {
		return nil, validationError
	}

	return &config, nil
}

func (c *Config) setDefaults() {
	defaultMotionPath, _ := exec.LookPath("motion")

	c.Port = 8888
	c.Host = ""
	c.MotionPath = defaultMotionPath
}

func (c *Config) validate() error {
	if c.Port == 0 {
		return errors.New("App port is invalid")
	}

	if c.MotionPath == "" {
		return errors.New("Path to motion is invalid or motion is not available")
	}

	return nil
}

func (c *Config) ListCamsToStart() []string {
	availableCams := c.Cameras
	toStart := []string{}

	for name, config := range availableCams {
		if config.Autostart == true {
			toStart = append(toStart, name)
		}
	}

	return toStart
}

func (c *Config) GetCameraOptions(cameraID string) (*CameraOptions, error) {
	availableCams := c.Cameras

	for id, config := range availableCams {
		if id == cameraID {
			return config, nil
		}
	}

	return nil, fmt.Errorf("No options available for camera %s", cameraID)
}
