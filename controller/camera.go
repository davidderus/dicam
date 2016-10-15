package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/davidderus/dicam/config"
)

type camera struct {
	id         string
	pid        int
	configFile string
	logFile    string
}

// ConfigDirectory is where the main and thread config are stored
const ConfigDirectory = "config"

// MainConfigFile is the default motion config
const MainConfigFile = "motion.conf"

// LogsDirectory is where the motion logs are stored
const LogsDirectory = "logs"

// ThreadBaseName is the model name for a thread configuration file
const ThreadBaseName = "dicam-thread-%s"

// DefaultConfigMode is the file mode for a config file
const DefaultConfigMode = 0644

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

	configError := c.buildConfig()
	if configError != nil {
		return configError
	}

	return nil
}

func (c *camera) buildConfig() error {
	mainConfigPath := path.Join(ConfigDirectory, MainConfigFile)
	defaultConfig, readError := ioutil.ReadFile(mainConfigPath)

	if readError != nil {
		return errors.New("Can not read main config file")
	}

	threadName := fmt.Sprintf(ThreadBaseName, c.id)
	c.configFile = path.Join(ConfigDirectory, threadName+".conf")
	c.logFile = path.Join(LogsDirectory, threadName+".log")

	// @note For now config is the default hard coded config
	configBytes := []byte(defaultConfig)
	writeError := ioutil.WriteFile(c.configFile, configBytes, DefaultConfigMode)

	if writeError != nil {
		return writeError
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
