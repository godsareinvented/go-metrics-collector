package main

import (
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/server"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
)

func main() {
	memStorage := mem_storage.NewInstance()
	repository.NewInstance(memStorage)

	webServer := server.Server{}
	webServer.Start()
}
