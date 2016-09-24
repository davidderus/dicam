package main

import "os/exec"

type camera struct {
  command *exec.Cmd
  configFile string
  running bool
}

func (c *camera) start() {
  command := exec.Command("motion")
  err := command.Start()

  if err != nil {
    panic(err)
  }

  c.command = command
  c.running = true

  return
}

func (c *camera) stop() {
  c.command.Process.Kill()
  c.running = false

  return
}

func main() {
  cam := camera{}
  cam.start()
}
