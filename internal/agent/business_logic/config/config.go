package config

import (
	"go.uber.org/zap"
	"time"
)

type Config struct {
	// Адрес эндпоинта HTTP-сервера
	Endpoint string `env:"ADDRESS" validate:"required,hostname_port"`

	// Частота отправки метрик на сервер в секундах
	// Значение поля должно быть >= значению поля PollInterval, иначе нет смысла отправлять не собранные метрики
	ReportInterval time.Duration `env:"REPORT_INTERVAL" validate:"required,gte=1,lte=86400,gtefield=PollInterval"`

	// Частота опроса метрик в секундах
	PollInterval time.Duration `env:"POLL_INTERVAL" validate:"required,gte=1,lte=86400"`

	// Ключ для вычисления хеша (hmac)
	HashKey string `env:"KEY" validate:"required,hashkey,min=3,max=512"`

	// Количество одновременно исходящих запросов на сервер (количество воркеров)
	// Значение поля должно быть не меньше чем ReportInterval / PollInterval, иначе приложение начнёт пропускать метрики
	RateLimit int `env:"RATE_LIMIT" validate:"required,ratelimmin,gte=1,lte=20"`

	// Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера
	GzipAcceptedContentTypes []string `validate:"required,unique,min=1,max=253,dive,oneof=application/json text/html"`

	// Минимальный размер тела ответа сервера, при котором будет происходить сжатие
	GzipMinContentLength int `validate:"required,gte=1,lt=10000"`

	// Логгер
	Logger *zap.Logger `validate:"required"`

	// Количество логических процессоров на текущей машине
	LogicalCpuCount int `validate:"required,numeric,gte=0,lte=1000"`

	// Список названий метрик, собираемый агентом. Список генерируемый
	MetricsToCollect []string `validate:"required,unique,min=1,max=1031,dive,mname"`
}

var Configuration Config
