package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/server"
)

func main() {
	webServer := server.Server{}
	webServer.Start()
}
