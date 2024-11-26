package config

import (
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"go.uber.org/zap"
)

type Config struct {
	Endpoint                 string                       `env:"ADDRESS"`           // Адрес эндпоинта HTTP-сервера.
	ReportInterval           int                          `env:"REPORT_INTERVAL"`   // Частота отправки метрик на сервер.
	PollInterval             int                          `env:"POLL_INTERVAL"`     // Частота опроса метрик из пакета runtime.
	StoreInterval            int                          `env:"STORE_INTERVAL"`    // Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск.
	FileStoragePath          string                       `env:"FILE_STORAGE_PATH"` // Путь до файла, куда сохраняются текущие значения.
	Restore                  bool                         `env:"RESTORE"`           // Булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера.
	DatabaseDSN              string                       `env:"DATABASE_DSN"`      // Адрес подключения к БД
	GzipAcceptedContentTypes []string                     // Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера
	GzipMinContentLength     int                          // Минимальный размер тела ответа сервера, при котором будет происходить сжатие
	PermanentStorage         *interfaces.PermanentStorage // Сконфигурированное постоянное хранилище метрик между работой сервера
	Repository               *repository.Repository       // Сконфигурированный репозиторий
	Logger                   *zap.Logger                  // Логгер
}

var Configuration Config
