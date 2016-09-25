package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

// Server logic

type commandServer struct {
	host string
	port int
}

func requestError(connection net.Conn, errorMessage string) {
	message := "ERROR"

	if errorMessage != "" {
		message = fmt.Sprintf("%s-%s", message, errorMessage)
	}

	connection.Write([]byte(message + "\n"))
	connection.Close()
}

func (cs *commandServer) start() error {
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
		requestError(connection, "")
	}

	parsedCommand := parseCommand(strings.TrimRight(string(message), "\n"))

	// todo Return and handle any error
	commandRunner(parsedCommand)

	connection.Close()
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
		default:
			return InvalidCommand{Command{}}
		}
	}
}

func commandRunner(command CommandInterface) {
	command.run()
}

// Command Handling

type CommandInterface interface {
	run() string
}

type Command struct{ params []string }

type CamCommand struct{ Command }

type ServerCommand struct{ Command }

type invalidCommand struct{ Command }

func (com CamCommand) run() (string, error) {
	return "", nil
}

func (com ServerCommand) run() (string, error) {
	return "", nil
}

func (com invalidCommand) run() (string, error) {
	return "", errors.New("Invalid command")
}
