package main

import (
	"github.com/oldhanasong/go-metrics-collector/internal/server"
)

func main() {
	webServer := server.Server{}
	webServer.Start()
}
