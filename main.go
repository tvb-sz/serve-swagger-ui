package main

import (
	"github.com/tvb-sz/serve-swagger-ui/app/console/command"
	_ "go.uber.org/automaxprocs"
)

func main() {
	command.Start()
}
