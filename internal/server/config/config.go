package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/server/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/server/repository"
	"go.uber.org/zap"
)

type Config struct {
	// Адрес эндпоинта HTTP-сервера
	Endpoint string `env:"ADDRESS" validate:"required,hostname_port"`

	// Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск
	StoreInterval int `env:"STORE_INTERVAL" validate:"required,numeric,gte=0,lte=86400"`

	// Путь до файла, куда сохраняются текущие значения
	FileStoragePath string `env:"FILE_STORAGE_PATH" validate:"required"`

	// Булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера
	Restore bool `env:"RESTORE" validate:"required"`

	// Адрес подключения к БД
	DatabaseDSN string `env:"DATABASE_DSN" validate:"omitempty"`

	// Ключ для вычисления хеша (hmac)
	HashKey string `env:"KEY" validate:"omitempty,hash_key,min=3,max=512"`

	// Разрешённые значения для заголовка "Content-Type" при сжатии ответа сервера
	GzipAcceptedContentTypes []string `validate:"required,unique,min=1,max=253,dive,oneof=application/json text/html"`

	// Минимальный размер тела ответа сервера, при котором будет происходить сжатие
	GzipMinContentLength int `validate:"required,numeric,gte=0,lt=10000"`

	// Сконфигурированное постоянное хранилище метрик между работой сервера
	PermanentStorage *interfaces.PermanentStorage `validate:"required"`

	// Сконфигурированный репозиторий
	Repository *repository.Repository `validate:"required"`

	// Логгер
	Logger *zap.Logger `validate:"required"`

	// Валидатор
	Validate *validator.Validate `validate:"required"`
}

var Configuration Config
