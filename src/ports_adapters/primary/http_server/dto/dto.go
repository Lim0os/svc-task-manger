package dto

import (
	"errors"
	"fmt"
	"svc-task_master/src/domain"
	"time"
)

// TaskRequest структура запроса для создания задачи
// swagger:model TaskRequest
type TaskRequest struct {
	// Тип задачи (обязательное поле)
	// required: true
	// example: "email_send"
	Type string `json:"type"`

	// Приоритет задачи
	// required: true
	// enum: low,medium,high,critical
	// example: "high"
	Priority string `json:"priority"`

	// Время планируемого выполнения
	// example: "2024-01-15T10:00:00Z"
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`

	// Полезная нагрузка задачи
	// required: true
	// example: {"email": "user@example.com", "subject": "Test"}
	Payload map[string]interface{} `json:"payload"`

	// Метаданные задачи
	// example: {"source": "api", "user_id": "123"}
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Количество попыток выполнения
	// minimum: 0
	// example: 0
	RetryCount int `json:"retryCount"`

	// Максимальное количество попыток
	// minimum: 0
	// example: 3
	MaxRetries int `json:"maxRetries"`

	// ID родительской задачи
	// example: "task-123"
	ParentTaskID string `json:"parentTaskId,omitempty"`

	// Список зависимостей
	// example: ["task-456", "task-789"]
	DependsOn []string `json:"dependsOn,omitempty"`

	// Очередь для выполнения
	// required: true
	// example: "default"
	Queue string `json:"queue"`
}

func (t *TaskRequest) Validate() error {

	if t.Type == "" {
		return errors.New("task type is required")
	}

	validPriorities := map[string]bool{
		string(domain.TaskPriorityLow):      true,
		string(domain.TaskPriorityMedium):   true,
		string(domain.TaskPriorityHigh):     true,
		string(domain.TaskPriorityCritical): true,
	}

	if t.Priority == "" {
		return errors.New("task priority is required")
	}
	if !validPriorities[t.Priority] {
		return fmt.Errorf("invalid priority: %s, must be one of: low, medium, high, critical", t.Priority)
	}

	if t.MaxRetries < 0 {
		return errors.New("max retries cannot be negative")
	}
	if t.RetryCount < 0 {
		return errors.New("retry count cannot be negative")
	}
	if t.RetryCount > t.MaxRetries {
		return errors.New("retry count cannot exceed max retries")
	}

	if t.ScheduledAt != nil && t.ScheduledAt.Before(time.Now()) {
		return errors.New("scheduled time must be in the future")
	}

	if t.Payload == nil {
		return errors.New("payload cannot be nil")
	}

	if t.Queue == "" {
		return errors.New("queue is required")
	}

	return nil
}

// GetTaskRequest структура запроса для получения задачи по ID
// swagger:model GetTaskRequest
type GetTaskRequest struct {
	// ID задачи
	// required: true
	// example: "task-123"
	ID string `json:"id"`
}

// UpdateTaskStatusRequest структура запроса для обновления статуса задачи
// swagger:model UpdateTaskStatusRequest
type UpdateTaskStatusRequest struct {
	// Новый статус задачи
	// required: true
	// enum: pending,processing,completed,failed,retrying
	// example: "completed"
	Status string `json:"status"`

	Id string `json:"id"`
}

func (t *UpdateTaskStatusRequest) Validate() error {
	if t.Status == "" {
		return errors.New("status is required")
	}
	if t.Id == "" {
		return errors.New("task id is required")
	}
	validStatuses := map[string]bool{
		string(domain.TaskStatusPending):    true,
		string(domain.TaskStatusProcessing): true,
		string(domain.TaskStatusCompleted):  true,
		string(domain.TaskStatusFailed):     true,
		string(domain.TaskStatusRetrying):   true,
	}

	if !validStatuses[t.Status] {
		return fmt.Errorf(
			"invalid status: %s, must be one of: pending, processing, completed, failed, retrying",
			t.Status,
		)
	}
	return nil
}

func (t *GetTaskRequest) Validate() error {
	if t.ID == "" {
		return errors.New("task id is required")
	}
	return nil
}

// GetTaskWhithFiltersRequest структура запроса для получения задач с фильтрами
// swagger:model GetTaskWhithFiltersRequest
type GetTaskWhithFiltersRequest struct {
	// Статус для фильтрации задач
	// enum: pending,processing,completed,failed,retrying
	// example: "pending"
	Status string `json:"status"`
}

func (r *GetTaskWhithFiltersRequest) Validate() error {
	if r.Status == "" {
		return nil
	}
	validStatuses := map[string]bool{
		string(domain.TaskStatusPending):    true,
		string(domain.TaskStatusProcessing): true,
		string(domain.TaskStatusCompleted):  true,
		string(domain.TaskStatusFailed):     true,
		string(domain.TaskStatusRetrying):   true,
	}

	if !validStatuses[r.Status] {
		return fmt.Errorf(
			"invalid status: %s, must be one of: pending, processing, completed, failed, retrying",
			r.Status,
		)
	}

	return nil
}

// Response универсальная структура ответа API
// swagger:model Response
type Response struct {
	// HTTP-статус код
	// example: 200
	Status int `json:"status,omitempty"`

	// Данные ответа (меняется в зависимости от endpoint)
	Data interface{} `json:"data,omitempty"`

	// Сообщение об ошибке (если есть)
	// example: "invalid request"
	Error *string `json:"error,omitempty"`
}
