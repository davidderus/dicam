package controller

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

// Server logic

const responseErrorCode = "ERROR"
const responseSuccessCode = "SUCCESS"

const invalidCommandError = "Invalid command"

const startAction = "START"
const stopAction = "STOP"
const listAction = "LIST"

// CommandCenter stores the CommandCenter listener config
type CommandCenter struct {
	Host string
	Port int
}

func sendResponse(connection net.Conn, responseType string, responseMessage string) {
	message := fmt.Sprintf("%s-%s", responseType, responseMessage)

	connection.Write([]byte(message + "\r"))
}

// Start starts the CommandCenter TCP listener
func (cs *CommandCenter) Start() error {
	service := fmt.Sprintf("%s:%d", cs.Host, cs.Port)
	tcpAddress, resolveError := net.ResolveTCPAddr("tcp4", service)
	if resolveError != nil {
		return resolveError
	}

	listener, listenError := net.Listen("tcp", tcpAddress.String())
	if listenError != nil {
		return listenError
	}

	defer listener.Close()

	for {
		connection, acceptError := listener.Accept()
		if acceptError != nil {
			return acceptError
		}

		go handleCommand(connection)
	}
}

// handleCommand receive a tcp command in a given format, parse the relevant
// informations and run the command before sending a response back
func handleCommand(connection net.Conn) {
	defer connection.Close()

	message, bufferError := bufio.NewReader(connection).ReadString('\r')

	if bufferError != nil {
		sendResponse(connection, responseErrorCode, bufferError.Error())
		return
	}

	trimmedMessage := strings.TrimRight(string(message), "\r")
	parsedCommand := parseCommand(trimmedMessage)

	output, runError := commandRunner(parsedCommand)

	var code, response string

	if runError != nil {
		code = responseErrorCode
		response = runError.Error()

		log.Printf("%s - %s - %s", code, trimmedMessage, response)
	} else {
		code = responseSuccessCode
		response = output

		log.Printf("%s - %s", code, trimmedMessage)
	}

	sendResponse(connection, code, response)
}

// Command Handling

type commandInterface interface {
	run() (string, error)
}

type command struct{ params []string }

type camCommand struct{ command }

type invalidCommand struct{ command }

// camCommand.run send a command against the CamsPoolInstance
// A camCommand must only provide a camera ID as an argument if needed
func (com camCommand) run() (string, error) {
	action := com.params[0]

	var id string

	if len(com.params) > 1 {
		id = com.params[1]
	} else {
		id = ""
	}

	switch action {
	case startAction:
		return CamsPoolInstance.launchCamera(id)
	case stopAction:
		return CamsPoolInstance.stopCamera(id)
	case listAction:
		return CamsPoolInstance.listCameras()
	}

	return "", errors.New(invalidCommandError)
}

func (com invalidCommand) run() (string, error) {
	return "", errors.New(invalidCommandError)
}

func commandRunner(command commandInterface) (string, error) {
	return command.run()
}

// parseCommand handles the command parsing.
// A command must respect the following format: "SUBJECT-COMMAND-ARG"
func parseCommand(input string) commandInterface {
	commandArray := strings.Split(input, "-")

	if len(commandArray) > 1 {
		mainCommand := commandArray[0]
		args := commandArray[1:]

		switch mainCommand {
		case "CAM":
			return camCommand{command{args}}
		}
	}

	return invalidCommand{command{}}
}
