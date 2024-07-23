package main

import (
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"github.com/godsareinvented/go-metrics-collector/internal/server"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/mem_storage"
)

func main() {
	memStorage := mem_storage.NewInstance()
	repository.NewInstance(memStorage)

	webServer := server.Server{}
	webServer.Start()
}
