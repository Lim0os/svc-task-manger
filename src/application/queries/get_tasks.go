package queries

import (
	"context"
	_ "errors"
	"svc-task_master/src/common/decorator"
	"svc-task_master/src/domain"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
)

type getTasksQuery struct {
	logger domain.ILogger
	repo   domain.IInMemoRepository
}

type GetTasksQuery decorator.CommandHandlerDecorator[dto.GetTaskWhithFiltersRequest, []domain.Task]

func NewGetTasksQuery(logger domain.ILogger, repo domain.IInMemoRepository) decorator.CommandHandlerDecorator[dto.GetTaskWhithFiltersRequest, []domain.Task] {
	return decorator.ApplyCommandLoggerDecorator[dto.GetTaskWhithFiltersRequest, []domain.Task](
		getTasksQuery{
			logger: logger,
			repo:   repo,
		},
		logger,
	)

}

func (c getTasksQuery) Handle(ctx context.Context, request dto.GetTaskWhithFiltersRequest) ([]domain.Task, error) {
	task, err := c.repo.GetAllFilterStatus(ctx, domain.TaskStatus(request.Status))
	return task, err
}
