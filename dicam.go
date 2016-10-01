package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/davidderus/dicam/controller"
	"github.com/urfave/cli"
)

const defaultPort = 8888

func main() {
	service := fmt.Sprintf("%s:%d", "", defaultPort)
	tcpAddress, resolveError := net.ResolveTCPAddr("tcp4", service)
	if resolveError != nil {
		log.Fatalln(resolveError.Error())
		return
	}

	sender, _ := net.Dial("tcp", tcpAddress.String())

	app := cli.NewApp()
	app.Name = "dicam-cli"
	app.Usage = "Controls dicam processes and cams"
	app.Version = dicamVersion

	app.Commands = []cli.Command{
		{
			Name:    "controller",
			Aliases: []string{"c"},
			Usage:   "Handles the main app control",
			Subcommands: []cli.Command{
				{
					Name:  "start",
					Usage: "Starts the controller",
					Action: func(c *cli.Context) error {
						log.Println("Starting command center")
						startError := controller.Start(defaultPort)

						if startError != nil {
							log.Fatalln(startError)
						}
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stops the controller",
					Action: func(c *cli.Context) error {
						return nil
					},
				},
			},
		},
		{
			Name:    "camera",
			Aliases: []string{"cam"},
			Usage:   "Interacts with a camera",
			Subcommands: []cli.Command{
				{
					Name:  "start",
					Usage: "Starts a camera",
					Action: func(c *cli.Context) error {
						fmt.Fprintf(sender, "CAM-START-1\n")

						response, _ := bufio.NewReader(sender).ReadString('\n')
						fmt.Println(response)
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stops a camera",
					Action: func(c *cli.Context) error {
						fmt.Fprintf(sender, "CAM-STOP-%s\n", c.Args().First())

						response, _ := bufio.NewReader(sender).ReadString('\n')
						fmt.Println(response)
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "Lists all available cameras",
					Action: func(c *cli.Context) error {
						fmt.Fprintf(sender, "CAM-LIST\n")

						response, _ := bufio.NewReader(sender).ReadString('\n')
						fmt.Println(response)
						return nil
					},
				},
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Manages the webserver",
			Subcommands: []cli.Command{
				{
					Name:  "start",
					Usage: "Starts the webserver",
					Action: func(c *cli.Context) error {
						fmt.Println("Starting webserver")
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stops the webserver",
					Action: func(c *cli.Context) error {
						fmt.Println("Stopping webserver")
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
