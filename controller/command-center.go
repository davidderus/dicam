package controller

import (
	"bufio"
	"errors"
	"fmt"
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
	Host       string
	Port       int
	controller *controller
}

func sendResponse(connection net.Conn, responseType string, responseMessage string) {
	message := fmt.Sprintf("%s-%s", responseType, responseMessage)

	connection.Write([]byte(message + "\n"))
	connection.Close()
}

func (cs *CommandCenter) Start() error {
	controller := &controller{}

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

		go handleCommand(connection, controller)
	}
}

func handleCommand(connection net.Conn, controller *controller) {
	message, bufferError := bufio.NewReader(connection).ReadString('\n')

	if bufferError != nil {
		sendResponse(connection, responseErrorCode, bufferError.Error())
	}

	parsedCommand := parseCommand(strings.TrimRight(string(message), "\n"))

	// todo Return and handle any error
	output, runError := commandRunner(parsedCommand, controller)

	if runError != nil {
		sendResponse(connection, responseErrorCode, runError.Error())
	} else {
		sendResponse(connection, responseSuccessCode, output)
	}

	connection.Close()
}

// Command Handling

type commandInterface interface {
	run(controller *controller) (string, error)
}

type command struct{ params []string }

type camCommand struct{ command }

type serverCommand struct{ command }

type invalidCommand struct{ command }

func (com camCommand) run(controller *controller) (string, error) {
	action := com.params[0]

	var id string

	if len(com.params) > 1 {
		id = com.params[1]
	} else {
		id = ""
	}

	switch action {
	case startAction:
		controller.launchCamera(1)
		return fmt.Sprintf("Starting cam %s", id), nil
	case stopAction:
		controller.stopCamera(1)
		return fmt.Sprintf("Stopping cam %s", id), nil
	case listAction:
		cams := controller.listCameras()
		var camsList []string

		for _, cam := range cams {
			camsList = append(camsList, fmt.Sprintf("Cam. %d - PID %d", cam.id, cam.pid))
		}

		return strings.Join(camsList, "\n"), nil
	}

	return "", errors.New(invalidCommandError)
}

func (com invalidCommand) run(controller *controller) (string, error) {
	return "", errors.New(invalidCommandError)
}

func commandRunner(command commandInterface, controller *controller) (string, error) {
	return command.run(controller)
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
