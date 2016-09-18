package main

import (
  "os"
  "fmt"
  "strings"
  "github.com/urfave/cli"
)

func main()  {
  cams := []string{"Cam 1", "Cam 2"}

  app := cli.NewApp()
  app.Name = "dicam-cli"
  app.Usage = "Controls dicam processes and cams"

  app.Commands = []cli.Command{
    {
      Name: "camera",
      Aliases: []string{"c"},
      Usage: "Interacts with a camera",
      Subcommands: []cli.Command{
        {
          Name: "start",
          Usage: "Starts a camera",
          Action: func(c *cli.Context) error {
            fmt.Println("Starting cam", c.Args().First())
            return nil
          },
        },
        {
          Name: "list",
          Usage: "Lists all available cameras",
          Action: func(c *cli.Context) error {
            fmt.Println(strings.Join(cams, "\n"))
            return nil
          },
        },
      },
    },
  }

  app.Run(os.Args)
}
