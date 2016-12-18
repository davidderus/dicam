package notifier

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

type EmailNotifier struct {
	Host     string
	Port     int
	From     string
	Password string
}

func (notifier *EmailNotifier) setOptions(options map[string]string) error {
	notifier.Host = options["host"]

	intPort, _ := strconv.Atoi(options["port"])
	notifier.Port = intPort

	notifier.From = options["from"]
	notifier.Password = options["password"]

	return nil
}

func (notifier *EmailNotifier) validateOptions() error {
	if notifier.Host == "" || notifier.Port == 0 {
		return errors.New("Invalid host or port in options")
	}

	if notifier.From == "" || notifier.Password == "" {
		return errors.New("A from email and an SMTP password are required")
	}

	return nil
}

func (notifier *EmailNotifier) send(message string, recipients []string) error {
	log.Printf("Sending email to %s from %s\n", strings.Join(recipients, ", "), notifier.From)

	smtpMessage := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		strings.Join(recipients, ","),
		"Email from dicam",
		message,
	)

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		notifier.From,
		notifier.Password,
		notifier.Host,
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", notifier.Host, notifier.Port),
		auth,
		notifier.From,
		recipients,
		[]byte(smtpMessage),
	)

	return err
}
