package notifier

import (
	"fmt"
	"strings"
)

type PushbulletNotifier struct {
	APIKey string
}

func (notifier PushbulletNotifier) send(message string, recipients []string, options map[string]string) error {
	fmt.Printf("Sending push to %s\n", strings.Join(recipients, ", "))
	return nil
}
