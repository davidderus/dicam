package notifier

import (
	"fmt"
	"strings"
)

type PushbulletNotifier struct {
	ApiKey string
}

func (notifier PushbulletNotifier) send(message string, recipients []string) error {
	fmt.Printf("Sending push to %s", strings.Join(recipients, ", "))
	return nil
}
