package logger

import (
	"log/slog"
	"os"
)

type CustomAsyncLogger struct {
	logger *slog.Logger
	ch     chan LogEntry
}

type LogEntry struct {
	Level   slog.Level
	Message string
	Attrs   []slog.Attr
}

func NewCustomAsyncLogger(bufferSize int, lvl string) *CustomAsyncLogger {
	l := &CustomAsyncLogger{
		logger: initLogger(lvl, false),
		ch:     make(chan LogEntry, bufferSize),
	}
	go l.startWorker()
	return l
}

func initLogger(lvl string, addSource bool) *slog.Logger {
	logLvl := parseLogLevel(lvl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       logLvl,
		AddSource:   addSource,
		ReplaceAttr: replaceAttr,
	}))
	return logger
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Value.Kind() == slog.KindDuration {
		return slog.String(attr.Key, attr.Value.Duration().String())
	}
	return attr
}

func (l *CustomAsyncLogger) startWorker() {
	for entry := range l.ch {

		args := make([]any, 0, len(entry.Attrs)*2)
		for _, attr := range entry.Attrs {
			args = append(args, attr.Key, attr.Value)
		}

		switch entry.Level {
		case slog.LevelDebug:
			l.logger.Debug(entry.Message, args...)
		case slog.LevelInfo:
			l.logger.Info(entry.Message, args...)
		case slog.LevelWarn:
			l.logger.Warn(entry.Message, args...)
		case slog.LevelError:
			l.logger.Error(entry.Message, args...)
		}
	}
}

func (l *CustomAsyncLogger) Info(msg string, attrs ...slog.Attr) {
	l.ch <- LogEntry{
		Level:   slog.LevelInfo,
		Message: msg,
		Attrs:   attrs,
	}
}

func (l *CustomAsyncLogger) Error(msg string, attrs ...slog.Attr) {
	l.ch <- LogEntry{
		Level:   slog.LevelError,
		Message: msg,
		Attrs:   attrs,
	}
}

func (l *CustomAsyncLogger) Debug(msg string, attrs ...slog.Attr) {
	l.ch <- LogEntry{
		Level:   slog.LevelDebug,
		Message: msg,
		Attrs:   attrs,
	}
}

func (l *CustomAsyncLogger) Warn(msg string, attrs ...slog.Attr) {
	l.ch <- LogEntry{
		Level:   slog.LevelWarn,
		Message: msg,
		Attrs:   attrs,
	}
}

func (l *CustomAsyncLogger) Shutdown() {
	close(l.ch)
}
