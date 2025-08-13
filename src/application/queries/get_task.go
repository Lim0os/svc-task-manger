package queries

import (
	"context"
	"errors"
	"svc-task_master/src/common/decorator"
	"svc-task_master/src/domain"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

type getTaskIdQuery struct {
	logger domain.ILogger
	repo   domain.IInMemoRepository
}

type GetTaskIdQuery decorator.CommandHandlerDecorator[dto.GetTaskRequest, domain.Task]

func NewGetTaskIdQuery(logger domain.ILogger, repo domain.IInMemoRepository) decorator.CommandHandlerDecorator[dto.GetTaskRequest, domain.Task] {
	return decorator.ApplyCommandLoggerDecorator[dto.GetTaskRequest, domain.Task](
		getTaskIdQuery{
			logger: logger,
			repo:   repo,
		},
		logger,
	)

}

func (c getTaskIdQuery) Handle(ctx context.Context, request dto.GetTaskRequest) (domain.Task, error) {
	task, ok := c.repo.Get(request.ID)
	if !ok {
		return domain.Task{}, errors.New("task not found")
	}
	return task, nil
}
