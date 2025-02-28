package config

import (
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Endpoint                 string        `env:"ADDRESS"`         // Адрес эндпоинта HTTP-сервера.
	ReportInterval           time.Duration `env:"REPORT_INTERVAL"` // Частота отправки метрик на сервер.
	PollInterval             time.Duration `env:"POLL_INTERVAL"`   // Частота опроса метрик из пакета runtime.
	HashKey                  string        `env:"KEY"`             // Ключ для вычисления хеша.
	GzipAcceptedContentTypes []string      // Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера
	GzipMinContentLength     int           // Минимальный размер тела ответа сервера, при котором будет происходить сжатие
	Logger                   *zap.Logger   // Логгер
	LogicalCpuCount          int           // Количество логических процессоров на текущей машине.
	MetricsToCollect         []string      // Список названий метрик, собираемый агентом. Список генерируемый.
}

var Configuration Config
