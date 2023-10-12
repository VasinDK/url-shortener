package slogdiscard

import (
	"context"

	"log/slog"
)

type DiscardHandler struct{}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Просто игнорируем запись журнала
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Возвращает тот же отбаботчик т.к. нет атребутов для сохранения
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Возвращает тот же отбаботчик т.к. нет группы для сохранения
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращается false, т.к. запись журнала игнорируется
	return false
}

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}
