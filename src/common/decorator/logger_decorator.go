package decorator

import (
	"context"
	"fmt"
	"log/slog"
	"svc-task_master/src/domain"
	"time"
)

type CommandLoggingDecorator[C any, R any] struct {
	base   CommandHandlerDecorator[C, R]
	logger domain.ILogger
}

func (d CommandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (R, error) {
	start := time.Now()
	handlerType := generateActionName(cmd)

	d.logger.Debug("Executing command",
		slog.String("command", handlerType),
		slog.String("command_body", fmt.Sprintf("%#v", cmd)),
	)
	result, err := d.base.Handle(ctx, cmd)
	defer func() {
		duration := time.Since(start)

		if err != nil {
			d.logger.Error("Failed to execute command",
				slog.String("command", handlerType),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()),
			)

		} else {
			d.logger.Info("Command executed successfully",
				slog.String("command", handlerType),
				slog.Duration("duration", duration),
			)
		}
	}()

	return result, err
}
