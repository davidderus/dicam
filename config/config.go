package config

import "github.com/spf13/viper"

type Config struct {
	Options *viper.Viper
}

func Read() (*Config, error) {
	options := viper.New()

	options.SetConfigName("config")
	options.AddConfigPath("$HOME/.config/dicam")
	err := options.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	config.Options = options

	return config, nil
}

func (c *Config) Cameras() []string {
	return c.Options.GetStringSlice("cameras")
}
