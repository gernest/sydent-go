// Package logger provides structured logging for various services or API
// endpoints for matrix exposed by matrixid.
//
// We are using bur/zap as the basis for our logger because of its performance
// and clean API.
package logger

import (
	"go.uber.org/zap"
)

func New() (Logger, error) {
	lg, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &zapper{log: lg}, nil
}

type Logger interface {
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

type zapper struct {
	log *zap.Logger
}

func (z *zapper) Error(msg string, fields ...zap.Field) {
	z.log.Error(msg, fields...)
}

func (z *zapper) Info(msg string, fields ...zap.Field) {
	z.log.Info(msg, fields...)
}

func (z *zapper) With(fields ...zap.Field) Logger {
	return &zapper{log: z.log.With(fields...)}
}

func (z *zapper) Sync() error {
	return z.log.Sync()
}
