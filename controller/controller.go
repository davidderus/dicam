package controller

import "github.com/davidderus/dicam/config"

var CamsPoolInstance *CamsPool

func Start(config *config.Config) error {
	CamsPoolInstance = &CamsPool{config: config}
	CamsPoolInstance.boot()

	cc := &CommandCenter{Port: config.Port}
	startError := cc.Start()

	if startError != nil {
		return startError
	}

	return nil
}
