package domain

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusRetrying   TaskStatus = "retrying"
)

type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

// Task представляет задачу в системе
// swagger:model Task
type Task struct {
	// Уникальный идентификатор задачи
	// example: "task-123"
	ID string `json:"id"`

	// Тип задачи
	// example: "email_send"
	Type string `json:"type"`

	// Текущий статус задачи
	// enum: pending,processing,completed,failed,retrying
	// example: "pending"
	Status TaskStatus `json:"status"`

	// Приоритет задачи
	// enum: low,medium,high,critical
	// example: "high"
	Priority TaskPriority `json:"priority"`

	// Время создания задачи
	// example: "2024-01-15T09:00:00Z"
	CreatedAt time.Time `json:"createdAt"`

	// Время последнего обновления
	// example: "2024-01-15T09:00:00Z"
	UpdatedAt time.Time `json:"updatedAt"`

	// Время планируемого выполнения
	// example: "2024-01-15T10:00:00Z"
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`

	// Время начала выполнения
	// example: "2024-01-15T10:00:00Z"
	StartedAt *time.Time `json:"startedAt,omitempty"`

	// Время завершения выполнения
	// example: "2024-01-15T10:05:00Z"
	FinishedAt *time.Time `json:"finishedAt,omitempty"`

	// Полезная нагрузка задачи
	// example: {"email": "user@example.com", "subject": "Test"}
	Payload map[string]interface{} `json:"payload"`

	// Метаданные задачи
	// example: {"source": "api", "user_id": "123"}
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Количество попыток выполнения
	// example: 0
	RetryCount int `json:"retryCount"`

	// Максимальное количество попыток
	// example: 3
	MaxRetries int `json:"maxRetries"`

	// Информация о последней ошибке
	LastError *TaskError `json:"lastError,omitempty"`

	// ID родительской задачи
	// example: "task-123"
	ParentTaskID string `json:"parentTaskId,omitempty"`

	// Список зависимостей
	// example: ["task-456", "task-789"]
	DependsOn []string `json:"dependsOn,omitempty"`

	// Очередь для выполнения
	// example: "default"
	Queue string `json:"queue,omitempty"`

	// ID воркера, выполняющего задачу
	// example: "worker-1"
	WorkerID string `json:"workerId,omitempty"`

	// Результат выполнения задачи
	Output interface{} `json:"output,omitempty"`
}

// TaskError представляет информацию об ошибке задачи
// swagger:model TaskError
type TaskError struct {
	// Сообщение об ошибке
	// example: "Failed to send email"
	Message string `json:"message"`

	// Стек вызовов (если доступен)
	// example: "main.main()\n\tmain.go:25 +0x123"
	Stack string `json:"stack,omitempty"`

	// Код ошибки
	// example: "EMAIL_SEND_FAILED"
	Code string `json:"code,omitempty"`
}
