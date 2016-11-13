// Package notifier handles the motion detection event given by motion.
// It wait a given time and then propagates the event to notifiers.
package notifier

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Event is a motion at a given time with an optionnal attachment
type Event struct {
	EventType string
	CameraID  string
	DateTime  time.Time
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
}

// Store logs the event in dicam database
func (e *Event) store() {
	println(e.EventType, "in", e.CameraID, "at", e.DateTime.Format(time.RFC1123))
}

// startCountdown waits for a given amount of seconds before sending a
// notification
func (e *Event) startCountdown() {
	fmt.Printf("Sending notification in %d seconds\n", defaultWaitTime)
	time.Sleep(defaultWaitTime * time.Second)

	if e.eventFile.filePath != "" {
		fmt.Printf("With one %s: %s\n", e.eventFile.fileType, e.eventFile.filePath)
	}

	e.notify()
}

// TODO Notify the user with a given string and file
func (e *Event) notify() {
	println("Notification sent!")
}

// convertFileType convers a motion fileTypeBit to an understandable one
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