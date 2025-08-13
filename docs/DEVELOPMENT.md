# Руководство разработчика Task Master

## Обзор архитектуры

Task Master построен по принципам Clean Architecture с использованием паттерна CQRS (Command Query Responsibility Segregation) и гексагональной архитектуры.

### Слои архитектуры

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP API (Primary Adapter)               │
├─────────────────────────────────────────────────────────────┤
│                  Application Layer (CQRS)                   │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │    Commands     │  │     Queries     │                  │
│  └─────────────────┘  └─────────────────┘                  │
├─────────────────────────────────────────────────────────────┤
│                    Domain Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │    Entities     │  │   Interfaces    │                  │
│  └─────────────────┘  └─────────────────┘                  │
├─────────────────────────────────────────────────────────────┤
│                  In-Memory DB (Secondary Adapter)          │
└─────────────────────────────────────────────────────────────┘
```

## Детальное описание компонентов

### 1. Domain Layer (`src/domain/`)

Содержит бизнес-логику и основные сущности системы.

#### Основные сущности

- **Task** - основная сущность задачи
- **TaskStatus** - enum статусов задачи
- **TaskPriority** - enum приоритетов задачи
- **TaskError** - информация об ошибках

#### Интерфейсы

- **TaskRepository** - интерфейс для работы с хранилищем задач

### 2. Application Layer (`src/application/`)

Реализует бизнес-логику приложения, разделяя команды и запросы.

#### Commands (`src/application/commands/`)

- **CreateTask** - создание новой задачи
- **UpdateTask** - обновление существующей задачи
- **Helper** - вспомогательные функции

#### Queries (`src/application/queries/`)

- **GetTask** - получение задачи по ID
- **GetTasks** - получение списка задач с фильтрацией

### 3. Ports & Adapters (`src/ports_adapters/`)

#### Primary Adapters (`src/ports_adapters/primary/`)

HTTP API сервер с REST endpoints:

- **Server** - основная структура сервера
- **Router** - маршрутизация HTTP запросов
- **DTO** - структуры для передачи данных
- **Handlers** - обработчики HTTP запросов

#### Secondary Adapters (`src/ports_adapters/secondary/`)

- **In-Memory DB** - высокопроизводительное хранилище
- **Service Layer** - сервисы приложения

### 4. Common (`src/common/`)

Общие утилиты и компоненты:

- **Config** - загрузка конфигурации из переменных окружения
- **Logger** - асинхронное логирование с батчингом
- **Decorator** - декораторы для команд и логирования

## Паттерны проектирования

### 1. CQRS (Command Query Responsibility Segregation)

Система разделяет операции на команды (изменяющие состояние) и запросы (читающие данные):

```go
// Command - изменяет состояние
type CreateTaskCommand struct {
    Type       string
    Priority   string
    Payload    map[string]interface{}
    // ... другие поля
}

// Query - читает данные
type GetTaskQuery struct {
    ID string
}
```

### 2. Repository Pattern

Абстракция доступа к данным через интерфейс:

```go
type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    GetByID(ctx context.Context, id string) (*Task, error)
    Update(ctx context.Context, task *Task) error
    GetByStatus(ctx context.Context, status TaskStatus) ([]*Task, error)
}
```

### 3. Decorator Pattern

Используется для добавления функциональности без изменения основного кода:

```go
type CommandHandler interface {
    Handle(ctx context.Context, command interface{}) error
}

type LoggingDecorator struct {
    handler CommandHandler
    logger  Logger
}
```

## Логирование

### Асинхронное логирование

Система использует асинхронное логирование с батчингом для высокой производительности:

```go
type CustomAsyncLogger struct {
    logger   *slog.Logger
    logChan  chan LogEntry
    batchSize int
    // ... другие поля
}

func (l *CustomAsyncLogger) Info(msg string, args ...any) {
    l.logChan <- LogEntry{
        Level: slog.LevelInfo,
        Message: msg,
        Args: args,
    }
}
```

### Уровни логирования

- **DEBUG** - детальная информация для разработки
- **INFO** - общая информация о работе системы
- **WARN** - предупреждения
- **ERROR** - ошибки

