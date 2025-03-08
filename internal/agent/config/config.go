package config

import "github.com/go-playground/validator/v10"

type Config struct {
	// Адрес эндпоинта HTTP-сервера
	Endpoint string `env:"ADDRESS" validate:"required,hostname_port"`

	// Частота отправки метрик на сервер в секундах
	// Значение поля должно быть >= значению поля PollInterval, т.к. нет смысла отправлять не собранные метрики
	ReportInterval int `env:"REPORT_INTERVAL" validate:"required,numeric,gte=0,lte=86400,gtefield=PollInterval"`

	// Частота опроса метрик в секундах
	PollInterval int `env:"POLL_INTERVAL" validate:"required,numeric,gte=0,lte=86400"`

	// Ключ для вычисления хеша (hmac)
	HashKey string `env:"KEY" validate:"required,hash_key,min=3,max=512"`

	// Количество одновременно исходящих запросов на сервер (количество воркеров)
	RateLimit int `env:"RATE_LIMIT" validate:"required,numeric,gte=0,lte=10"`

	// Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера
	GzipAcceptedContentTypes []string `validate:"required,unique,min=1,max=253,dive,oneof=application/json text/html"`

	// Минимальный размер тела ответа сервера, при котором будет происходить сжатие
	GzipMinContentLength int `validate:"required,numeric,gte=0,lt=10000"`

	// Колиество логических процессоров на текущей машине
	LogicalCpuCount int `validate:"required,numeric,gte=0,lte=1000"`

	// Список названий метрик, собираемый агентом. Список генерируемый
	MetricNameList []string `validate:"required,unique,min=1,max=1031,dive,metric_name"`

	// Валидатор
	Validate *validator.Validate `validate:"required"`
}

var Configuration Config
