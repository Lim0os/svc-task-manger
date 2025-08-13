package commands

import (
	"context"
	"svc-task_master/src/common/decorator"
	"svc-task_master/src/domain"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

type createTaskCommnad struct {
	logger domain.ILogger
	repo   domain.IInMemoRepository
}

type CreateTaskCommnad decorator.CommandHandlerDecorator[dto.TaskRequest, string]

func NewCreateTaskCommnad(logger domain.ILogger, repo domain.IInMemoRepository) decorator.CommandHandlerDecorator[dto.TaskRequest, string] {
	return decorator.ApplyCommandLoggerDecorator[dto.TaskRequest, string](
		createTaskCommnad{
			logger: logger,
			repo:   repo,
		},
		logger,
	)

}

func (c createTaskCommnad) Handle(ctx context.Context, request dto.TaskRequest) (string, error) {
	task := createTask(request)
	c.repo.SetUpdate(task.ID, task)
	return task.ID, nil
}
