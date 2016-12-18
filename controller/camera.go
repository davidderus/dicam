package controller

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"

	"github.com/davidderus/dicam/config"
)

// camera is a physical camera declared in the configuration file with an
// optionnal PID if it is started
type camera struct {
	ID           string
	pid          int
	StreamPort   int
	configFile   string
	logFile      string
	workingDir   string
	NotifierPath string
	UserOptions  *config.CameraOptions
}

func (c *camera) setWorkingDir(directory string) {
	c.workingDir = directory
}

func (c *camera) setStreamPort(streamPort int) {
	c.StreamPort = streamPort
}

// getNotifierPath gets the binary current path in order to launch watchers
func (c *camera) getNotifierPath() error {
	appDir, appDirError := filepath.Abs(filepath.Dir(os.Args[0]))
	if appDirError != nil {
		return appDirError
	}

	c.NotifierPath = appDir

	return nil
}

// setup initiates build of a camera config file for motion after checking
// some basic items
func (c *camera) setup(cameraOptions *config.CameraOptions) error {
	if len(c.ID) == 0 {
		return errors.New("No id set for camera")
	}

	if len(cameraOptions.Device) == 0 {
		return errors.New("No device infos set")
	}

	_, deviceStatError := os.Stat(cameraOptions.Device)
	if deviceStatError != nil {
		if os.IsNotExist(deviceStatError) {
			return fmt.Errorf("Device %s not found, aborting.", cameraOptions.Device)
		}

		return deviceStatError
	}

	c.UserOptions = cameraOptions

	watcherError := c.getNotifierPath()
	if watcherError != nil {
		return watcherError
	}

	configError := c.buildConfig()
	if configError != nil {
		return configError
	}

	return nil
}

// buildConfig builds a motion config file from scratch with the user-defined
// camera options.
//
// TODO Do not rewrite config file if options are unchanged
func (c *camera) buildConfig() error {
	mainConfigPath := path.Join(config.TemplatesDirectory, config.MainConfigFileTemplate)

	threadName := fmt.Sprintf(config.ThreadBaseName, c.ID)
	c.configFile = path.Join(c.workingDir, config.ConfigDirectoryName, threadName+".conf")
	c.logFile = path.Join(c.workingDir, config.LogsDirectoryName, threadName+".log")

	// Read from default template
	template, parseError := template.ParseFiles(mainConfigPath)
	if parseError != nil {
		return errors.New("Can not read nor parse main config template: " + parseError.Error())
	}

	// Execute config options against template
	outputConfig, configError := os.Create(c.configFile)
	if configError != nil {
		return errors.New("Can not open thread config")
	}

	defer outputConfig.Close()

	// Write to file
	templateExecuteError := template.Execute(outputConfig, c)
	if templateExecuteError != nil {
		return templateExecuteError
	}

	return nil
}

// start launches a motion binary for the camera
func (c *camera) start() error {
	command := exec.Command("motion", "-c", c.configFile, "-l", c.logFile)
	err := command.Start()

	if err != nil {
		return err
	}

	c.pid = command.Process.Pid

	command.Process.Release()

	return nil
}

func (c *camera) stop() error {
	err := syscall.Kill(c.pid, syscall.SIGTERM)

	if err != nil {
		return err
	}

	c.pid = 0

	return nil
}

func (c *camera) isRunning() bool {
	return c.pid > 0
}

func (c camera) teardown() {
	os.Remove(c.configFile)
}

func (c *camera) infos() string {
	infosTemplate := "DEVICE:%s\nROLE:%s\nPID:%d\nPORT:%d\nLOG:%s\nCONFIG:%s"

	return fmt.Sprintf(
		infosTemplate,
		c.ID,
		c.UserOptions.Role,
		c.pid,
		c.StreamPort,
		c.logFile,
		c.configFile,
	)
}
