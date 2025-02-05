package config

type Config struct {
	Endpoint                 string   `env:"ADDRESS"`         // Адрес эндпоинта HTTP-сервера.
	ReportInterval           int      `env:"REPORT_INTERVAL"` // Частота отправки метрик на сервер.
	PollInterval             int      `env:"POLL_INTERVAL"`   // Частота опроса метрик из пакета runtime.
	GzipAcceptedContentTypes []string // Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера
	GzipMinContentLength     int      // Минимальный размер тела ответа сервера, при котором будет происходить сжатие
}

var Configuration Config
