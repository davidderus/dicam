package watcher

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Event is a motion at a given time
type Event struct {
	EventType string
	CameraID  string
	DateTime  time.Time
	eventFile EventFile
}

// EventFile is a file linked to an event
type EventFile struct {
	filePath string
	fileType string
}

const defaultWaitTime = 10

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
func (e *Event) AddFile(filePath string, fileTypeBit int) {
	fileType, convertError := convertFileType(fileTypeBit)
	if convertError != nil {
		log.Fatalln(convertError)
	}

	e.eventFile = EventFile{filePath, fileType}
}

func (e *Event) Trigger() {
	e.store()
	e.startCountdown()
}

// Store logs the event in dicam database
func (e *Event) store() {
	println(e.EventType, "in", e.CameraID, "at", e.DateTime.Format(time.RFC1123))
}

// TODO Wait a given time before alerting the end user
func (e *Event) startCountdown() {
	e.notify()
}

// TODO Notify the user with a given string and file
func (e *Event) notify() {
	fmt.Printf("Sending notification in %d seconds\n", defaultWaitTime)

	if e.eventFile.filePath != "" {
		fmt.Printf("With one %s: %s\n", e.eventFile.fileType, e.eventFile.filePath)
	}

	time.Sleep(defaultWaitTime * time.Second)

	println("Notification sent!")
}

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
