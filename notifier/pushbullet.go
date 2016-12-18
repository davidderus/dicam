package notifier

import (
	"errors"
	"log"
	"strings"

	pushbullet "github.com/xconstruct/go-pushbullet"
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
	pusher := pushbullet.New(notifier.APIKey)

	log.Printf("Sending push to %s\n", strings.Join(recipients, ", "))

	pushErrors := []string{}

	for _, recipient := range recipients {
		pushError := pusher.PushNote(recipient, "Push from dicam", message)

		if pushError != nil {
			pushErrors = append(pushErrors, pushError.Error())
		}
	}

	if len(pushErrors) > 0 {
		errorStrings := strings.Join(pushErrors, ", ")
		return errors.New("Errors when pushing to recipients: " + errorStrings)
	}

	return nil
}
