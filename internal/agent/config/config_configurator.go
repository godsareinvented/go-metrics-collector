package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type ConfigConfigurator struct{}

func (c *ConfigConfigurator) ParseConfig() {
	Configuration = Config{
		GzipAcceptedContentTypes: []string{"application/json", "text/html"},
		GzipMinContentLength:     1400,
	}

	parseFlags()
	err := parseEnv()
	if nil != err {
		panic("Error parsing env: " + err.Error())
	}
}

func parseFlags() {
	flag.StringVar(&Configuration.Endpoint, "a", ":8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&Configuration.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&Configuration.PollInterval, "p", 2, "Частота опроса метрик из пакета runtime")

	flag.Parse()
	// todo: Отрицательные значения?
}

func parseEnv() error {
	return env.Parse(&Configuration)
}
