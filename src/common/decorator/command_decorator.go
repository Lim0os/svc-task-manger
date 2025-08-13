package decorator

import (
	"context"
	"fmt"

	"strings"
	"svc-task_master/src/domain"
)

type CommandHandlerDecorator[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

type SemaDecorator[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

func ApplyCommandLoggerDecorator[C any, R any](
	handler CommandHandlerDecorator[C, R],
	logger domain.ILogger,
) CommandHandlerDecorator[C, R] {
	return CommandLoggingDecorator[C, R]{
		logger: logger,
		base:   handler,
	}
}
