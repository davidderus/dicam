// Package notifier handles the motion detection event given by motion.
// It wait a given time and then propagates the event to notifiers.
package notifier

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/davidderus/dicam/config"
)

// Event is a motion at a given time with an optionnal attachment
type Event struct {
	EventType string
	CameraID  string
	DateTime  time.Time
	Config    *config.Config
	eventFile EventFile
}

// EventFile is an attachment linked to an event
type EventFile struct {
	filePath string
	fileType string
}

// defaultWaitTime is the time before firing an event
//
// This is set in order not to immediately alert when detecting a motion and
// letting some time for the user to deactivate the notifier (ie: when entering
// his property)
const defaultWaitTime = 10

// SetDateTime parse the epoch time given by motion and update the Event with
// the normalized value
func (e *Event) SetDateTime(motionTime string) error {
	unixTime, parseError := strconv.ParseInt(motionTime, 10, 64)
	if parseError != nil {
		log.Fatalf("%s is not a valid epoch time", motionTime)
	}

	parsedDateTime := time.Unix(unixTime, 0)
	if parsedDateTime.IsZero() {
		log.Fatalln("Can't parse given time")
	}

	e.DateTime = parsedDateTime

	return nil
}

// AddFile adds files to the current Event
// See convertFileType for allowed fileTypeBit
func (e *Event) AddFile(filePath string, fileTypeBit int) {
	fileType, convertError := convertFileType(fileTypeBit)
	if convertError != nil {
		log.Fatalln(convertError)
	}

	e.eventFile = EventFile{filePath, fileType}
}

// Trigger stores the event infos and starts the countdown
func (e *Event) Trigger() {
	e.store()
	e.startCountdown()
	e.notify()
}

// Store logs the event in dicam database
func (e *Event) store() {
	log.Println(e.EventType, "in", e.CameraID, "at", e.DateTime.Format(time.RFC1123))
}

// startCountdown waits for a given amount of seconds before sending a
// notification
func (e *Event) startCountdown() {
	waitTime := e.Config.Countdown

	if waitTime == 0 {
		waitTime = defaultWaitTime
	}

	log.Printf("Sending notification in %d seconds\n", waitTime)
	time.Sleep(time.Duration(waitTime) * time.Second)

	if e.eventFile.filePath != "" {
		log.Printf("With one %s: %s\n", e.eventFile.fileType, e.eventFile.filePath)
	}
}

// TODO Notify the user with a given string and file
func (e *Event) notify() {
	if len(e.Config.Notifiers) == 0 {
		log.Fatalf("No notifiers in config, aborting")
	}

	for _, notifierConfig := range e.Config.Notifiers {
		var notifier notifierInterface

		// Getting notifier
		notifier, optionsError := getNotifier(notifierConfig.Service, notifierConfig.ServiceOptions)

		if optionsError != nil {
			// We do not use a Fatalf which could prevent execution of other notifiers
			log.Printf("%s: %s", notifierConfig.Service, optionsError.Error())
			continue
		}

		// Sending notification
		notifyError := notifier.send("azerty", notifierConfig.Recipients)

		if notifyError != nil {
			// We do not use a Fatalf which could prevent execution of other notifiers
			log.Printf("%s: %s\n", notifierConfig.Service, notifyError.Error())
			continue
		}

		log.Printf("Notification sent to %s recipients!\n", notifierConfig.Service)
	}
}

type notifierInterface interface {
	// send sends a message to a given set of recipients
	send(message string, recipients []string) error

	// setOptions defines configuration options for the notifier
	setOptions(options map[string]string) error

	// validateOptions validates the options given by setOptions and required by
	// the notifier
	validateOptions() error
}

type invalidNotifier struct{}

func (notifier *invalidNotifier) send(message string, recipients []string) error {
	return errors.New("Invalid notifier in config, no notification sent.")
}

func (notifier *invalidNotifier) setOptions(options map[string]string) error {
	return nil
}

func (notifier *invalidNotifier) validateOptions() error {
	return nil
}

func getNotifier(service string, options map[string]string) (notifierInterface, error) {
	var notifier notifierInterface

	switch service {
	case "pushbullet", "push":
		notifier = &PushbulletNotifier{}
	case "email", "mail":
		notifier = &EmailNotifier{}
	default:
		return &invalidNotifier{}, nil
	}

	notifierOptionsError := notifier.setOptions(options)

	if notifierOptionsError != nil {
		return notifier, notifierOptionsError
	}

	notifierValidationError := notifier.validateOptions()

	return notifier, notifierValidationError
}

// convertFileType converts a motion fileTypeBit to an understandable one
func convertFileType(fileTypeBit int) (string, error) {
	knownFormats := map[int]string{
		1:  "Normal picture",
		2:  "Snapshot picture",
		4:  "Debug picture",
		8:  "Movie file",
		16: "Debug Movie File",
		32: "Timelapse",
	}

	format, exists := knownFormats[fileTypeBit]
	if exists {
		return format, nil
	}

	return "", errors.New("Unknown motion filetype")
}
