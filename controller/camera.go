package controller

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/davidderus/dicam/config"
)

// camera is a physical camera declared in the configuration file with an
// optionnal PID if it is started
type camera struct {
	ID           string
	StreamPort   int
	CapturesDir  string
	NotifierPath string
	UserOptions  *config.CameraOptions

	pid int

	workingDir string
	configFile string
	logFile    string
}

func (c *camera) setWorkingDir(directory string) {
	c.workingDir = directory
}

func (c *camera) setStreamPort(streamPort int) {
	c.StreamPort = streamPort
}

// getNotifierPath gets the absolute path to the dicam binary in order to launch watchers
func (c *camera) getNotifierPath() error {
	appDir, appDirError := filepath.Abs(os.Args[0])
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

	if cameraOptions.GetCameraType() == "" {
		return errors.New("Invalid camera type. You must have one of `device` (local device address) or `remote_device` (remote camera IP)")
	}

	cameraType := cameraOptions.GetCameraType()

	var pingError error

	// Checking local or remote camera existence
	if cameraType == "local" {
		pingError = pingLocalCamera(cameraOptions.Device)
	} else {
		pingError = pingRemoteCamera(cameraOptions.RemoteDevice, cameraOptions.RemoteDeviceAuth)
	}

	if pingError != nil {
		return pingError
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
	c.CapturesDir = path.Join(c.workingDir, config.CapturesDirectoryName, threadName)

	// Read from default template
	templateData, assetError := Asset(mainConfigPath)
	if assetError != nil {
		return errors.New("Can not read nor parse main config template: " + assetError.Error())
	}

	template, parseError := template.New("motion.conf").Parse(string(templateData))
	if parseError != nil {
		return errors.New("Can not read nor parse main config template: " + parseError.Error())
	}

	// Ensure that captures directory exists
	mkdirCamCapturesDir := os.MkdirAll(c.CapturesDir, 0700)
	if mkdirCamCapturesDir != nil {
		return errors.New("Can not create captures storage directory")
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

// pingLocalCamera check the existence of the camera on the host.
func pingLocalCamera(deviceAddress string) error {
	_, deviceStatError := os.Stat(deviceAddress)
	if deviceStatError != nil {
		if os.IsNotExist(deviceStatError) {
			return fmt.Errorf("Device %s not found, aborting", deviceAddress)
		}

		return deviceStatError
	}

	return nil
}

// pingRemoteCamera issue a GET request to the remote camera URL, but does not read the body.
// If the server respond with a 200, then the camera is accepted.
// We're not doing HEAD as some webcam server may return a body on HEAD resulting
// in a `http.ErrBodyNotAllowed`
func pingRemoteCamera(remoteAddress, remoteAuth string) error {
	httpClient := &http.Client{}
	httpRequest, _ := http.NewRequest("GET", remoteAddress, nil)

	if remoteAuth != "" {
		remoteAuthArray := strings.Split(remoteAuth, ":")
		httpRequest.SetBasicAuth(remoteAuthArray[0], remoteAuthArray[1])
	}

	response, headError := httpClient.Do(httpRequest)
	response.Body.Close()

	if headError != nil {
		return headError
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Expecting a HTTP 200 for %s, but %d found", remoteAddress, response.StatusCode)
	}

	return nil
}
