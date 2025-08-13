package application

import (
	"svc-task_master/src/application/commands"
	"svc-task_master/src/application/queries"
)

type App struct {
	Command Commands
	Query   Queries
}

type Commands struct {
	CreateTask commands.CreateTaskCommnad
	UpdateTask commands.UpdateTaskCommnad
}

type Queries struct {
	GetTask  queries.GetTaskIdQuery
	GetTasks queries.GetTasksQuery
}
