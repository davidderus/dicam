package main

import (
	"fmt"
	"net"
)

type commandServer struct {
	host string
	port int
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

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	buffer := make([]byte, 1024)
	connection.Read(buffer)

	connection.Write([]byte("hello from dicam\n"))
	connection.Close()
}
