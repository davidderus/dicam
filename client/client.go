package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

// Client allows communication with the CommandCenter
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

// Ask sends a request to the CommandCenter and returns its response
func (c *Client) Ask(command string) (string, error) {
	fmt.Fprintf(c.sender, command+"\r")

	output, _ := bufio.NewReader(c.sender).ReadString('\r')
	output = strings.TrimRight(string(output), "\r")

	response := strings.SplitN(output, "-", 2)

	if len(response) > 1 {
		if response[0] != "SUCCESS" {
			return fmt.Sprintf("%s: %s", response[0], response[1]), nil
		} else {
			return "", errors.New(fmt.Sprintln(response[1]))
		}
	} else {
		return "", errors.New(fmt.Sprintln("Unknown response from command center"))
	}
}

// Print logs the response of the command center
// TODO Use a real cli logger to differentiate error from success
func (c *Client) Print(command string) {
	askResponse, askError := c.Ask(command)

	if askError != nil {
		fmt.Println(askError)
	} else {
		fmt.Println(askResponse)
	}
}
