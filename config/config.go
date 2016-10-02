package config

import (
	"errors"
	"os/exec"

	"github.com/spf13/viper"
)

type Camera struct {
	Device    string
	Role      string
	Autostart bool `mapstructure:"auto_start"`
	Notifiers []*Notifier
	Watcher   *Watcher
}

type Notifier struct {
	Service    string
	Recipients []string
}

type Watcher struct {
	AutoStart bool `mapstructure:"auto_start"`
	Countdown int
}

type Config struct {
	Port       int
	Host       string
	MotionPath string `mapstructure:"motion_path"`
	Cameras    map[string]*Camera
}

func Read() (*Config, error) {
	options := viper.New()

	options.SetConfigName("config")
	options.AddConfigPath("$HOME/.config/dicam")
	setDefaultOptions(options)

	readError := options.ReadInConfig()
	if readError != nil {
		return nil, readError
	}

	validationError := validateOptions(options)
	if validationError != nil {
		return nil, validationError
	}

	var config Config
	unmarshalError := options.Unmarshal(&config)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	return &config, nil
}

func setDefaultOptions(options *viper.Viper) {
	defaultMotionPath, _ := exec.LookPath("motion")

	// Setting default for non nested values
	// See https://github.com/spf13/viper/issues/162 for related issue
	options.SetDefault("port", 8888)
	options.SetDefault("host", "")
	options.SetDefault("motion_path", defaultMotionPath)
}

func validateOptions(options *viper.Viper) error {
	if !options.IsSet("port") || options.GetInt("port") == 0 {
		return errors.New("App port is invalid")
	}

	if options.Get("motionPath") == "" {
		return errors.New("Path to motion is invalid or motion is not available")
	}

	return nil
}
