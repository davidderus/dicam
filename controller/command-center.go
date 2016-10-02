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

type CommandCenter struct {
	Host string
	Port int
}

func sendResponse(connection net.Conn, responseType string, responseMessage string) {
	message := fmt.Sprintf("%s-%s", responseType, responseMessage)

	connection.Write([]byte(message + "\n"))
}

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

func handleCommand(connection net.Conn) {
	defer connection.Close()

	message, bufferError := bufio.NewReader(connection).ReadString('\n')

	if bufferError != nil {
		sendResponse(connection, responseErrorCode, bufferError.Error())
		return
	}

	parsedCommand := parseCommand(strings.TrimRight(string(message), "\n"))

	output, runError := commandRunner(parsedCommand)

	if runError != nil {
		sendResponse(connection, responseErrorCode, runError.Error())
		log.Println(runError)
	} else {
		sendResponse(connection, responseSuccessCode, output)
		log.Println(output)
	}
}

// Command Handling

type commandInterface interface {
	run() (string, error)
}

type command struct{ params []string }

type camCommand struct{ command }

type serverCommand struct{ command }

type invalidCommand struct{ command }

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
