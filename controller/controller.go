package controller

var CamsPoolInstance *CamsPool

func Start(port int) error {
	CamsPoolInstance = &CamsPool{}

	cc := &CommandCenter{Port: port}
	startError := cc.Start()

	if startError != nil {
		return startError
	}

	return nil
}
