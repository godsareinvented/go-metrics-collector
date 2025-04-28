package main

import (
	"context"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/general/utils"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/server"
	"github.com/godsareinvented/go-metrics-collector/internal/server/server/callback"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exitCh := make(chan os.Signal, 1)
		signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
		<-exitCh
		cancel()
	}()

	configConfigurator := config.ConfigConfigurator{}
	err := configConfigurator.Parse(ctx)
	if nil != err {
		panic(err)
	}

	webServer := server.Server{
		OnStart: callback.OnServerStartedCallback,
		OnStop:  callback.OnServerStoppedCallback,
	}
	wg, errCh := webServer.Start(ctx, cancel)

	// todo: Необходимо переписать схему. Прекрасно работает для одиночных ошибок.
	//  Но появись функция, генерирующая ошибки постоянно, при возникновении которых нельзя останавливать сервер, -
	//  будет переполняться переменная err.
	//  Необходимо решение с разными логгерами и некоторым хранилищем логеров.
	go func() {
		for {
			select {
			case errNew, ok := <-errCh:
				if nil != errNew {
					err = utils.WrapErrs(err, errNew)
				}
				if !ok {
					err = utils.WrapErrs(errors.New("error channel has been closed"), err)
					panic(err)
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		if nil != wg {
			wg.Wait()
		}
		panic(errors.New("context canceled"))
	}
}
