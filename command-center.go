package main

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
	host string
	port int
}

func sendResponse(connection net.Conn, responseType string, responseMessage string) {
	message := fmt.Sprintf("%s-%s", responseType, responseMessage)

	connection.Write([]byte(message + "\n"))
	connection.Close()
}

func (cs *CommandCenter) start() error {
	service := fmt.Sprintf("%s:%d", cs.host, cs.port)
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

type CommandInterface interface {
	run() (string, error)
}

type Command struct{ params []string }

type CamCommand struct{ Command }

type ServerCommand struct{ Command }

type InvalidCommand struct{ Command }

func (com CamCommand) run() (string, error) {
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

func (com ServerCommand) run() (string, error) {
	action := com.params[0]

	switch action {
	case startAction:
		return "Starting webserver", nil
	case stopAction:
		return "Stopping webserver", nil
	}

	return "", errors.New(invalidCommandError)
}

func (com InvalidCommand) run() (string, error) {
	return "", errors.New(invalidCommandError)
}

func commandRunner(command CommandInterface) (string, error) {
	return command.run()
}

func parseCommand(command string) CommandInterface {
	commandArray := strings.Split(command, "-")

	if len(commandArray) > 1 {
		command := commandArray[0]
		args := commandArray[1:]

		switch command {
		case "CAM":
			return CamCommand{Command{args}}
		case "SERVER":
			return ServerCommand{Command{args}}
		}
	}

	return InvalidCommand{Command{}}
}
