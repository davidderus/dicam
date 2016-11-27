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

func (notifier EmailNotifier) send(message string, recipients []string, options map[string]string) error {
	fmt.Printf("Sending email to %s from %s\n", strings.Join(recipients, ", "), options["from"])
	return nil
}
