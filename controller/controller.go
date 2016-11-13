package controller

import "github.com/davidderus/dicam/config"

// CamsPoolInstance allows global access to the CamsPool in the controller package
var CamsPoolInstance *CamsPool

// Start initialiazes a CamsPoolInstance and let the CommandCenter start
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
