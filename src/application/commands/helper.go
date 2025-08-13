package commands

import (
	"github.com/google/uuid"
	"svc-task_master/src/domain"
	"svc-task_master/src/ports_adapters/primary/http_server/dto"
	"time"
)

func createTask(request dto.TaskRequest) domain.Task {
	now := time.Now()
	id := uuid.New().String()

	return domain.Task{
		ID:           id,
		Type:         request.Type,
		Status:       domain.TaskStatusPending,
		Priority:     domain.TaskPriority(request.Priority),
		CreatedAt:    now,
		UpdatedAt:    now,
		ScheduledAt:  request.ScheduledAt,
		Payload:      request.Payload,
		Metadata:     request.Metadata,
		RetryCount:   request.RetryCount,
		MaxRetries:   request.MaxRetries,
		ParentTaskID: request.ParentTaskID,
		DependsOn:    request.DependsOn,
		Queue:        request.Queue,

		StartedAt:  nil,
		FinishedAt: nil,
		LastError:  nil,
		WorkerID:   "",
		Output:     nil,
	}
}
