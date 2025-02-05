package config

import (
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/server/repository"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Endpoint                 string                       `env:"ADDRESS"`         // Адрес эндпоинта HTTP-сервера.
	ReportInterval           time.Duration                `env:"REPORT_INTERVAL"` // Частота отправки метрик на сервер.
	PollInterval             time.Duration                `env:"POLL_INTERVAL"`   // Частота опроса метрик из пакета runtime.
	GzipAcceptedContentTypes []string                     // Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера
	GzipMinContentLength     int                          // Минимальный размер тела ответа сервера, при котором будет происходить сжатие
	PermanentStorage         *interfaces.PermanentStorage // Сконфигурированное постоянное хранилище метрик между работой сервера
	Repository               *repository.Repository       // Сконфигурированный репозиторий
	Logger                   *zap.Logger                  // Логгер
}

var Configuration Config
