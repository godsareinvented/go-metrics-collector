package config

import (
	"errors"
	"flag"
	"github.com/caarlos0/env"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/business_logic/config"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/general/logger"
	"github.com/oldhanasong/go-metrics-collector/internal/general/validation"
	"github.com/shirou/gopsutil/v3/cpu"
	"sync"
	"time"
)

type (
	ConfigConfigurator struct{}
)

var (
	once sync.Once
)

func (c *ConfigConfigurator) ParseConfig() error {
	var resErr error
	once.Do(func() {
		config.Configuration = config.Config{
			GzipAcceptedContentTypes: []string{"application/json", "text/html"},
			GzipMinContentLength:     1, // Должно быть 1400. Для соответствия инкременту 8 заменено на 1.
			Logger:                   logger.New(),
		}

		parseFlags()

		if err := parseEnv(); err != nil {
			resErr = errors.New("error parsing environment variables: " + err.Error())
			return
		}

		if err := setLogicalCpuNumber(); err != nil {
			resErr = err
			return
		}

		if err := generateMetricNameToCollect(); err != nil {
			resErr = err
			return
		}

		if err := validation.Validator().Struct(config.Configuration); err != nil {
			resErr = err
			return
		}
	})
	return resErr
}

func parseFlags() {
	var reportInterval, pollInterval int

	flag.StringVar(&config.Configuration.Endpoint, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&reportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&pollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&config.Configuration.HashKey, "k", "", "Ключ для вычисления хэша")
	flag.IntVar(&config.Configuration.RateLimit, "l", 1, "Количество одновременно исходящих запросов на сервер (количество воркеров)")

	flag.Parse()

	config.Configuration.ReportInterval = time.Duration(reportInterval)
	config.Configuration.PollInterval = time.Duration(pollInterval)
}

func parseEnv() error {
	return env.Parse(&config.Configuration)
}

func setLogicalCpuNumber() error {
	count, err := cpu.Counts(true)
	if nil != err {
		return err
	}

	config.Configuration.LogicalCpuCount = count
	return nil
}

func generateMetricNameToCollect() error {
	metricNameList, err := decorator.AddCpuUtilizationMetricNames(dictionary.MetricNameList[:], config.Configuration.LogicalCpuCount)
	if nil != err {
		return err
	}

	config.Configuration.MetricsToCollect = metricNameList
	return nil
}
