package cli

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Client defines basic client options
type Client struct {
	Host   string
	Port   int
	sender net.Conn
}

// Connect opens a tcp channel to the CommandCenter
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

// Ask sends a request to the CommandCenter
func (c *Client) Ask(command string) {
	fmt.Fprintf(c.sender, command+"\r")

	output, _ := bufio.NewReader(c.sender).ReadString('\r')
	output = strings.TrimRight(string(output), "\r")

	response := strings.SplitN(output, "-", 2)

	if len(response) > 1 {
		if response[0] != "SUCCESS" {
			fmt.Printf("%s: %s", response[0], response[1])
		} else {
			fmt.Println(response[1])
		}
	} else {
		fmt.Println("Unknown response from command center")
	}
}
