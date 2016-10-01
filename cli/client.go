package cli

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	Host   string
	Port   int
	sender net.Conn
}

func (c *Client) Connect() error {
	service := fmt.Sprintf("%s:%d", c.Host, c.Port)
	tcpAddress, resolveError := net.ResolveTCPAddr("tcp4", service)
	if resolveError != nil {
		return resolveError
	}

	sender, dialError := net.Dial("tcp", tcpAddress.String())

	if dialError != nil {
		return dialError
	}

	c.sender = sender

	return nil
}

func (c *Client) Ask(command string) string {
	fmt.Fprintf(c.sender, command+"\n")

	response, _ := bufio.NewReader(c.sender).ReadString('\n')

	return response
}
