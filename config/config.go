// Package config defines and parses configuration for controllers, cameras,
// notifiers and watchers
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

// CameraOptions lists the options allowed for a camera
type CameraOptions struct {
	// Device address like /dev/video0
	Device string

	// Basic Motion options
	Width     int
	Height    int
	Framerate int

	// MotionThreshold is `threshold` in Motion config
	MotionThreshold int `toml:"motion_threshold"`

	// EventGap is `gap` in Motion config
	EventGap int `toml:"event_gap"`

	// Role is one of []string{"stream", "watch"}
	Role string

	// Autostart defines if the camera should be started at boot
	Autostart bool `toml:"auto_start"`
}

// NotifierOptions includes the option for a notifier
type NotifierOptions struct {
	// Notifying service name
	Service string

	// Notifications recipients
	Recipients []string

	// Service options
	ServiceOptions map[string]string
}

// Config is the default config object
type Config struct {
	Port int
	Host string

	// Countdown before a notification is sent
	Countdown int

	// Path to motion binary
	MotionPath string `toml:"motion_path"`

	// Directory where logs and generated config files are stored
	WorkingDir string `toml:"working_dir"`

	// Listing of Camera with their options
	Cameras map[string]*CameraOptions

	// All cameras with a watch role will use the given Notifiers
	Notifiers map[string]*NotifierOptions
}

// TemplatesDirectory is where the main and thread config are stored
const TemplatesDirectory = "templates"

// ConfigDirectoryName is the name for the thread configs directory
const ConfigDirectoryName = "configs"

// LogsDirectoryName is the name for the directory where the motion logs are stored
const LogsDirectoryName = "logs"

// MainConfigFileTemplate is the default motion config
const MainConfigFileTemplate = "motion.conf.tpl"

// ThreadBaseName is the model name for a thread configuration file
const ThreadBaseName = "dicam-thread-%s"

// DefaultConfigMode is the file mode for a config file
const DefaultConfigMode = 0700

// Read reads config for dicam
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

	config.setDefaults(userHomeDir)

	populateError := config.populateWorkingDir()
	if populateError != nil {
		return nil, populateError
	}

	return &config, nil
}

// setDefaults defines default options in configuration such as motion path,
// controller port and host…
func (c *Config) setDefaults(userDir string) {
	defaultMotionPath, _ := exec.LookPath("motion")

	c.Port = 8888
	c.Host = ""
	c.MotionPath = defaultMotionPath
	c.WorkingDir = path.Join(userDir, ".dicam")
}

// validate validates a few config options to prevent further errors
func (c *Config) validate() error {
	if c.Port == 0 {
		return errors.New("App port is invalid")
	}

	if c.MotionPath == "" {
		return errors.New("Path to motion is invalid or motion is not available")
	}

	return nil
}

// populateWorkingDir creates the configs and logs directories based on the
// WorkingDir
func (c *Config) populateWorkingDir() error {
	userDirError := os.MkdirAll(c.WorkingDir, DefaultConfigMode)
	if userDirError != nil {
		return userDirError
	}

	mkdirConfigError := os.MkdirAll(path.Join(c.WorkingDir, ConfigDirectoryName), DefaultConfigMode)
	if mkdirConfigError != nil {
		return mkdirConfigError
	}

	mkdirLogsError := os.MkdirAll(path.Join(c.WorkingDir, LogsDirectoryName), DefaultConfigMode)
	if mkdirLogsError != nil {
		return mkdirLogsError
	}

	return nil
}

// ListCamsToStart returns ids of cameras to start at boot time
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

// GetCameraOptions returns the CameraOptions for a given cameraID
func (c *Config) GetCameraOptions(cameraID string) (*CameraOptions, error) {
	availableCams := c.Cameras

	for id, options := range availableCams {
		if id == cameraID {
			return options, nil
		}
	}

	return nil, fmt.Errorf("No options available for camera %s", cameraID)
}
