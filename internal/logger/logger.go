package logger

import "go.uber.org/zap"

func NewInstance() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()

	return logger
}
