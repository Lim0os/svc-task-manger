package domain

import "log/slog"

type ILogger interface {
	Info(msg string, attrs ...slog.Attr)
	Error(msg string, attrs ...slog.Attr)
	Debug(msg string, attrs ...slog.Attr)
	Warn(msg string, attrs ...slog.Attr)
}
