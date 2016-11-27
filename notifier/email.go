package notifier

import (
	"fmt"
	"strings"
)

type EmailNotifier struct {
	Host     string
	Port     int
	From     string
	Password string
}

func (notifier EmailNotifier) send(message string, recipients []string) error {
	fmt.Printf("Sending email to %s", strings.Join(recipients, ", "))
	return nil
}
