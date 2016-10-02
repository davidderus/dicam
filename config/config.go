package config

import (
	"errors"
	"os/exec"

	"github.com/spf13/viper"
)

type Config struct {
	Options *viper.Viper
}

func Read() (*Config, error) {
	options := viper.New()

	options.SetConfigName("config")
	options.AddConfigPath("$HOME/.config/dicam")
	setDefaultOptions(options)

	err := options.ReadInConfig()
	if err != nil {
		return nil, err
	}

	validationError := validateOptions(options)
	if validationError != nil {
		return nil, validationError
	}

	config := &Config{}
	config.Options = options

	return config, nil
}

func setDefaultOptions(options *viper.Viper) {
	defaultMotionPath, _ := exec.LookPath("motion")

	options.SetDefault("controller.port", 8888)
	options.SetDefault("controller.host", "")
	options.SetDefault("motionPath", defaultMotionPath)
}

func validateOptions(options *viper.Viper) error {
	if !options.IsSet("controller.port") || options.GetInt("controller.port") == 0 {
		return errors.New("Controller port is invalid")
	}

	if options.Get("motionPath") == "" {
		return errors.New("Path to motion is invalid or motion is not available")
	}

	return nil
}

func (c *Config) Cameras() map[string]interface{} {
	return c.Options.Sub("cameras").AllSettings()
}
