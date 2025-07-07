package config

import (
	"errors"
	"flag"
	"github.com/caarlos0/env"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/server/logger"
	"github.com/oldhanasong/go-metrics-collector/internal/server/service/validator"
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
		Configuration = Config{
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

		if err := generateMetricNameList(); err != nil {
			resErr = err
			return
		}

		if err := validator.GetValidator().Struct(Configuration); err != nil {
			resErr = err
			return
		}
	})
	return resErr
}

func parseFlags() {
	var reportInterval, pollInterval int

	flag.StringVar(&Configuration.Endpoint, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&reportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&pollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&Configuration.HashKey, "k", "", "Ключ для вычисления хэша")
	flag.IntVar(&Configuration.RateLimit, "l", 1, "Количество одновременно исходящих запросов на сервер (количество воркеров)")

	flag.Parse()

	Configuration.ReportInterval = time.Duration(reportInterval)
	Configuration.PollInterval = time.Duration(pollInterval)
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

	Configuration.MetricsToCollect = metricNameList
	return nil
}
