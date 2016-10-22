package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/davidderus/dicam/config"
	"github.com/davidderus/dicam/controller"
	"github.com/davidderus/dicam/watcher"

	"github.com/urfave/cli"
)

func getClient(config *config.Config) *Client {
	client := &Client{Host: config.Host, Port: config.Port}
	connectionError := client.Connect()

	if connectionError != nil {
		log.Fatalln(connectionError)
	}

	return client
}

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
				return nil
			},
		},
		{
			Name:   "watcher",
			Hidden: true,
			Action: func(c *cli.Context) error {
				eventType := c.Args().Get(1)
				watcherEvent := watcher.Event{
					CameraID:  c.Args().Get(0),
					EventType: eventType,
				}

				watcherEvent.SetDateTime(c.Args().Get(2))

				if eventType == "picture" {
					watcherEvent.AddFile(c.Args().Get(3), c.Args().Get(4))
				}

				watcherEvent.Store()

				return nil
			},
		},
	}

	app.Run(os.Args)
}
