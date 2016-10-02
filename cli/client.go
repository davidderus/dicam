package cli

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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

func (c *Client) Ask(command string) {
	fmt.Fprintf(c.sender, command+"\n")

	output, _ := bufio.NewReader(c.sender).ReadString('\n')
	output = strings.TrimRight(string(output), "\n")

	response := strings.SplitN(output, "-", 2)

	if len(response) > 1 {
		fmt.Printf("%s: %s", response[0], response[1])
	} else {
		fmt.Println("Unknown response from command center")
	}
}
