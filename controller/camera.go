package controller

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/davidderus/dicam/config"
)

type camera struct {
	id          string
	pid         int
	configFile  string
	logFile     string
	workingDir  string
	userOptions *config.CameraOptions
}

// MainConfigFileTemplate is the default motion config
const MainConfigFileTemplate = "motion.conf.tpl"

// ThreadBaseName is the model name for a thread configuration file
const ThreadBaseName = "dicam-thread-%s"

// DefaultConfigMode is the file mode for a config file
const DefaultConfigMode = 0644

func (c *camera) setWorkingDir(directory string) {
	c.workingDir = directory
}

func (c *camera) setup(cameraOptions *config.CameraOptions) error {
	if len(c.id) == 0 {
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

	c.userOptions = cameraOptions

	configError := c.buildConfig()
	if configError != nil {
		return configError
	}

	return nil
}

// todo: Do not rewrite config file if options are unchanged
func (c *camera) buildConfig() error {
	mainConfigPath := path.Join(config.TemplatesDirectoryName, MainConfigFileTemplate)

	threadName := fmt.Sprintf(ThreadBaseName, c.id)
	c.configFile = path.Join(c.workingDir, config.TemplatesDirectoryName, threadName+".conf")
	c.logFile = path.Join(c.workingDir, config.LogsDirectoryName, threadName+".log")

	// Read from default template
	template, parseError := template.ParseFiles(mainConfigPath)
	if parseError != nil {
		return errors.New("Can not read nor parse main config template")
	}

	// Execute config options against template
	outputConfig, configError := os.Create(c.configFile)
	if configError != nil {
		return errors.New("Can not open thread config")
	}

	defer outputConfig.Close()

	// Write to file
	templateExecuteError := template.Execute(outputConfig, c.userOptions)
	if templateExecuteError != nil {
		return templateExecuteError
	}

	return nil
}

func (c *camera) start() error {
	command := exec.Command("motion", "-c", c.configFile, "-l", c.logFile)
	err := command.Start()

	if err != nil {
		return err
	}

	c.pid = command.Process.Pid

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

func (c camera) isRunning() bool {
	return c.pid > 0
}

func (c camera) teardown() {
	os.Remove(c.configFile)
}
