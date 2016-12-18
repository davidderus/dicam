// Package cli defines all cli options and instanciates a client to communicate
// with the CommandCenter
package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/davidderus/dicam/config"
	"github.com/davidderus/dicam/controller"
	"github.com/davidderus/dicam/notifier"
	"github.com/davidderus/dicam/server"

	"github.com/urfave/cli"
)

// getClient creates a new client to send command to the CommandCenter
func getClient(config *config.Config) *Client {
	client := &Client{Host: config.Host, Port: config.Port}
	connectionError := client.Connect()

	if connectionError != nil {
		log.Fatalln(connectionError)
	}

	return client
}

// loadConfig reads the dicam config file
func loadConfig() *config.Config {
	config, readError := config.Read()
	if readError != nil {
		log.Fatalln(readError)
	}

	return config
}

// Init starts Dicam command line interface
func Init(version string) {
	var client *Client

	appConfig := loadConfig()

	app := cli.NewApp()
	app.Name = "dicam-cli"
	app.Usage = "Controls dicam processes and cams"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:    "controller",
			Aliases: []string{"c"},
			Usage:   "Starts the app control",
			Action: func(c *cli.Context) error {
				log.Printf("Starting command center on %s:%d", appConfig.Host, appConfig.Port)
				startError := controller.Start(appConfig)

				if startError != nil {
					log.Fatalln(startError)
				}

				return nil
			},
		},
		{
			Name:    "camera",
			Aliases: []string{"cam"},
			Usage:   "Interacts with a camera",
			Before: func(c *cli.Context) error {
				client = getClient(appConfig)
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "start",
					Usage: "Starts a camera",
					Action: func(c *cli.Context) error {
						client.Ask("CAM-START-" + c.Args().First())
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stops a camera",
					Action: func(c *cli.Context) error {
						client.Ask("CAM-STOP-" + c.Args().First())
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "Lists all available cameras",
					Action: func(c *cli.Context) error {
						client.Ask("CAM-LIST")
						return nil
					},
				},
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Starts the webserver",
			Action: func(c *cli.Context) error {
				fmt.Println("Starting webserver")
				server.Start()
				return nil
			},
		},
		{
			Name:   "notifier",
			Hidden: true,
			Action: func(c *cli.Context) error {
				cameraID := c.Args().Get(0)
				eventType := c.Args().Get(1)

				notifierEvent := notifier.Event{
					CameraID:  cameraID,
					EventType: eventType,
					Config:    appConfig,
				}

				epochTime := c.Args().Get(2)
				notifierEvent.SetDateTime(epochTime)

				if eventType == "pictureSave" {
					filePath := c.Args().Get(3)
					fileTypeBit, _ := strconv.Atoi(c.Args().Get(4))
					notifierEvent.AddFile(filePath, fileTypeBit)
				}

				notifierEvent.Trigger()

				return nil
			},
		},
	}

	app.Run(os.Args)
}
