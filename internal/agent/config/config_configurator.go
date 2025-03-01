package config

import (
	"errors"
	"flag"
	"github.com/caarlos0/env"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary/decorator"
	"github.com/shirou/gopsutil/v3/cpu"
)

type ConfigConfigurator struct{}

func (c *ConfigConfigurator) ParseConfig() error {
	Configuration = Config{
		GzipAcceptedContentTypes: []string{"application/json", "text/html"},
		GzipMinContentLength:     1400,
	}

	parseFlags()
	err := parseEnv()
	if nil != err {
		return errors.New("Error parsing environment variables: " + err.Error())
	}
	err = setLogicalCpuNumber()
	if nil != err {
		return err
	}
	err = generateMetricNameList()
	if nil != err {
		return err
	}
	return nil
}

func parseFlags() {
	flag.StringVar(&Configuration.Endpoint, "a", ":8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&Configuration.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&Configuration.PollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&Configuration.HashKey, "k", "", "Ключ для вычисления хэша")
	flag.IntVar(&Configuration.RateLimit, "l", 1, "Количество одновременно исходящих запросов на сервер (количество воркеров)")

	flag.Parse()
	// todo: Отрицательные значения?
}

func parseEnv() error {
	return env.Parse(&Configuration)
}

func setLogicalCpuNumber() error {
	count, err := cpu.Counts(true)
	if nil != err {
		return err
	}

	Configuration.LogicalCpuCount = count
	return nil
}

func generateMetricNameList() error {
	metricNameList, err := decorator.AddCpuUtilizationMetricNames(dictionary.MetricNameList[:], Configuration.LogicalCpuCount)
	if nil != err {
		return err
	}

	Configuration.MetricNameList = metricNameList
	return nil
}
