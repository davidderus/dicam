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
	Host string
	Port int
}

func sendResponse(connection net.Conn, responseType string, responseMessage string) {
	message := fmt.Sprintf("%s-%s", responseType, responseMessage)

	connection.Write([]byte(message + "\n"))
	connection.Close()
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
	message, bufferError := bufio.NewReader(connection).ReadString('\n')

	if bufferError != nil {
		sendResponse(connection, responseErrorCode, bufferError.Error())
	}

	parsedCommand := parseCommand(strings.TrimRight(string(message), "\n"))

	// todo Return and handle any error
	output, runError := commandRunner(parsedCommand)

	if runError != nil {
		sendResponse(connection, responseErrorCode, runError.Error())
	} else {
		sendResponse(connection, responseSuccessCode, output)
	}

	connection.Close()
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
		return fmt.Sprintf("Starting cam %s", id), nil
	case stopAction:
		return fmt.Sprintf("Stopping cam %s", id), nil
	case listAction:
		return "Listing all cams", nil
	}

	return "", errors.New(invalidCommandError)
}

func (com serverCommand) run() (string, error) {
	action := com.params[0]

	switch action {
	case startAction:
		return "Starting webserver", nil
	case stopAction:
		return "Stopping webserver", nil
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
		case "SERVER":
			return serverCommand{command{args}}
		}
	}

	return invalidCommand{command{}}
}
