package config

type Config struct {
	Endpoint                 string   `env:"ADDRESS"`         // Адрес эндпоинта HTTP-сервера.
	ReportInterval           int      `env:"REPORT_INTERVAL"` // Частота отправки метрик на сервер.
	PollInterval             int      `env:"POLL_INTERVAL"`   // Частота опроса метрик из пакета runtime.
	HashKey                  string   `env:"KEY"`             // Ключ для вычисления хеша.
	RateLimit                int      `env:"RATE_LIMIT"`      // Количество одновременно исходящих запросов на сервер (количество воркеров).
	GzipAcceptedContentTypes []string // Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера.
	GzipMinContentLength     int      // Минимальный размер тела ответа сервера, при котором будет происходить сжатие.
	LogicalCpuCount          int      // Колиество логических процессоров на текущей машине.
	MetricNameList           []string // Список названий метрик, собираемый агентом. Список генерируемый.
}

var Configuration Config
