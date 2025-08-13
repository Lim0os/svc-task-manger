package commands

import (
	"context"
	"svc-task_master/src/common/decorator"
	"svc-task_master/src/domain"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

type updateTaskCommnad struct {
	logger domain.ILogger
	repo   domain.IInMemoRepository
}

type UpdateTaskCommnad decorator.CommandHandlerDecorator[dto.UpdateTaskStatusRequest, any]

func NewUpdateTaskCommnad(logger domain.ILogger, repo domain.IInMemoRepository) decorator.CommandHandlerDecorator[dto.UpdateTaskStatusRequest, any] {
	return decorator.ApplyCommandLoggerDecorator[dto.UpdateTaskStatusRequest, any](
		updateTaskCommnad{
			logger: logger,
			repo:   repo,
		},
		logger,
	)

}

func (c updateTaskCommnad) Handle(ctx context.Context, request dto.UpdateTaskStatusRequest) (any, error) {
	c.repo.UpdateStatus(request.Id, domain.TaskStatus(request.Status))
	return nil, nil
}
