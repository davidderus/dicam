package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"syscall"
)

type camera struct {
	id         int
	pid        int
	configFile string
	logFile    string
	running    bool
}

// ConfigDirectory is where the main and thread config are stored
const ConfigDirectory = "config"

// MainConfigFile is the default motion config
const MainConfigFile = "motion.conf"

// LogsDirectory is where the motion logs are stored
const LogsDirectory = "logs"

// ThreadBaseName is the model name for a thread configuration file
const ThreadBaseName = "dicam-thread-%d"

// DefaultConfigMode is the file mode for a config file
const DefaultConfigMode = 0644

func (c *camera) setup() {
	if c.id == 0 {
		panic("No id set for camera")
	}

	mainConfigPath := path.Join(ConfigDirectory, MainConfigFile)
	defaultConfig, readError := ioutil.ReadFile(mainConfigPath)

	if readError != nil {
		panic(readError)
	}

	threadName := fmt.Sprintf(ThreadBaseName, c.id)
	c.configFile = path.Join(ConfigDirectory, threadName+".conf")
	c.logFile = path.Join(LogsDirectory, threadName+".log")

	// @note For now config is the default hard coded config
	configBytes := []byte(defaultConfig)
	writeError := ioutil.WriteFile(c.configFile, configBytes, DefaultConfigMode)

	if writeError != nil {
		panic(writeError)
	}
}

func (c *camera) start() {
	command := exec.Command("motion", "-c", c.configFile, "-l", c.logFile)
	err := command.Start()

	if err != nil {
		panic(err)
	}

	c.pid = command.Process.Pid
}

func (c *camera) stop() {
	err := syscall.Kill(c.pid, syscall.SIGTERM)

	if err != nil {
		panic(err)
	}

	c.pid = 0
}

func (c camera) isRunning() bool {
	return c.pid > 0
}

func (c camera) teardown() {
	os.Remove(c.configFile)
}
