package notifier

import (
	"errors"
	"fmt"
	"strings"
)

type PushbulletNotifier struct {
	APIKey string
}

func (notifier *PushbulletNotifier) setOptions(options map[string]string) error {
	notifier.APIKey = options["api_key"]
	return nil
}

func (notifier *PushbulletNotifier) validateOptions() error {
	if notifier.APIKey == "" {
		return errors.New("An API Key is needed")
	}

	return nil
}

func (notifier *PushbulletNotifier) send(message string, recipients []string) error {
	fmt.Printf("Sending push to %s\n", strings.Join(recipients, ", "))
	return nil
}
