package config

import (
	"errors"
	"fmt"
	"os"
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
	WorkingDir string `toml:"working_dir"`
	Cameras    map[string]*CameraOptions
}

// TemplatesDirectoryName is where the main and thread config are stored
const TemplatesDirectoryName = "templates"

const ConfigDirectoryName = "configs"

// LogsDirectoryName is where the motion logs are stored
const LogsDirectoryName = "logs"

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

	populateError := config.populateWorkingDir()
	if populateError != nil {
		return nil, populateError
	}

	return &config, nil
}

func (c *Config) setDefaults() {
	defaultMotionPath, _ := exec.LookPath("motion")

	c.Port = 8888
	c.Host = ""
	c.MotionPath = defaultMotionPath
	c.WorkingDir = path.Join(os.TempDir(), "dicam")
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

func (c *Config) populateWorkingDir() error {
	// Adds logs and config files directories
	mkdirConfigError := os.MkdirAll(path.Join(c.WorkingDir, ConfigDirectoryName), 0644)
	if mkdirConfigError != nil {
		return mkdirConfigError
	}

	mkdirLogsError := os.MkdirAll(path.Join(c.WorkingDir, LogsDirectoryName), 0644)
	if mkdirLogsError != nil {
		return mkdirLogsError
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

	for id, options := range availableCams {
		if id == cameraID {
			return options, nil
		}
	}

	return nil, fmt.Errorf("No options available for camera %s", cameraID)
}
