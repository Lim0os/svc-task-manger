package application

import (
	"svc-task_master/src/application"
	"svc-task_master/src/application/commands"
	"svc-task_master/src/application/queries"
	"svc-task_master/src/domain"
)

func InitApp(repo domain.IInMemoRepository, logger domain.ILogger) application.App {
	return application.App{
		Command: application.Commands{
			CreateTask: commands.NewCreateTaskCommnad(logger, repo),
			UpdateTask: commands.NewUpdateTaskCommnad(logger, repo),
		},
		Query: application.Queries{
			GetTasks: queries.NewGetTasksQuery(logger, repo),
			GetTask:  queries.NewGetTaskIdQuery(logger, repo),
		},
	}
}
