package db

import (
	"svc-task_master/src/domain"
	"svc-task_master/src/ports_adapters/secondary/inmemory/db/task_repo"
)

type Repository struct {
	InMemoryDB domain.IInMemoRepository
}

func NewRepository(logger domain.ILogger, sharedNum int) *Repository {
	return &Repository{
		InMemoryDB: task_repo.NewSharderStorage(sharedNum),
	}
}
