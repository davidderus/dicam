package notifier

import (
	"fmt"
	"strconv"
	"strings"
)

type EmailNotifier struct {
	Host     string
	Port     int
	From     string
	Password string
}

func (notifier EmailNotifier) setOptions(options map[string]string) error {
	notifier.Host = options["host"]

	intPort, _ := strconv.Atoi(options["port"])
	notifier.Port = intPort

	notifier.From = options["from"]
	notifier.Password = options["password"]

	return nil
}

func (notifier EmailNotifier) send(message string, recipients []string) error {
	fmt.Printf("Sending email to %s from %s", strings.Join(recipients, ", "), notifier.From)
	return nil
}
