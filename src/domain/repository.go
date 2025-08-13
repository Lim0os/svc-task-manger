package domain

import "context"

type IInMemoRepository interface {
	Get(key string) (Task, bool)
	SetUpdate(key string, data Task)
	GetAllFilterStatus(ctx context.Context, status TaskStatus) ([]Task, error)
	UpdateStatus(key string, status TaskStatus)
}
