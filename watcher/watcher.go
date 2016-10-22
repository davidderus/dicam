package watcher

import (
	"log"
	"strconv"
	"time"
)

// Event is a motion at a given time
type Event struct {
	EventType  string
	CameraID   string
	DateTime   time.Time
	eventFiles []EventFile
}

// EventFile is a file linked to an event
type EventFile struct {
	filePath string
	fileType string
}

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

// Addfile adds files to the current Event
func (e *Event) AddFile(filePath string, fileType string) {
	e.eventFiles = append(e.eventFiles, EventFile{filePath, fileType})
}

// Store logs the event in dicam database
func (e *Event) Store() {
	println(e.EventType, e.CameraID, e.DateTime.Format(time.RFC1123))
}

// TODO Wait a given time before alerting the end user
func (e *Event) startCountdown() {
}

// TODO Notify the user with a given string
func (e *Event) notify(withImage bool) {
}
