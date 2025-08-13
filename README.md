# Task Master - Система управления задачами

Task Master - это высокопроизводительная система управления задачами, построенная на архитектуре Clean Architecture с использованием Go. Система предоставляет REST API для создания, управления и мониторинга задач с поддержкой различных статусов, приоритетов и очередей.


## 🏗️ Архитектура

Проект построен по принципам Clean Architecture с четким разделением на слои:

```
src/
├── domain/           # Бизнес-логика и сущности
├── application/      # Слой приложения (команды и запросы)
├── ports_adapters/  # Адаптеры для внешних интерфейсов
│   ├── primary/     # HTTP API
│   └── secondary/   # In-memory хранилище
└── common/          # Общие утилиты (конфигурация, логирование)
```

### Основные компоненты

- **Domain Layer** - содержит бизнес-сущности (Task, TaskStatus, TaskPriority)
- **Application Layer** - реализует команды (CreateTask, UpdateTask) и запросы (GetTask, GetTasks)
- **Primary Adapters** - HTTP сервер с REST API endpoints
- **Secondary Adapters** - In-memory хранилище с шардированием
- **Common** - конфигурация, логирование, декораторы

## Требования

- Go 1.23.0 или выше
- Docker (опционально, для контейнеризации)


## Запуск приложения
```bash
go run main.go
```

### Docker

```bash
# Сборка образа
docker build -t task-master .

# Запуск контейнера
docker run -p 8080:8080 task-master
```

## ⚙️ Конфигурация

Система настраивается через переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `PORT` | Порт HTTP сервера | `8080` |
| `LOG_LEVEL` | Уровень логирования | `debug` |
| `BATCH_SIZE` | Размер батча для логирования | `100` |
| `MEMORY_TTL` | TTL для in-memory данных (сек) | `300` |
| `NUM_SHARDS` | Количество шардов для БД | `100` |

### Пример .env файла
```env
PORT=8080
LOG_LEVEL=info
BATCH_SIZE=50
MEMORY_TTL=600
NUM_SHARDS=50
```

## 📚 API Endpoints

### Создание задачи
```http
POST /task
Content-Type: application/json

{
  "type": "email_send",
  "priority": "high",
  "payload": {
    "email": "user@example.com",
    "subject": "Test Email"
  },
  "queue": "default",
  "scheduledAt": "2024-01-15T10:00:00Z"
}
```

### Получение задачи по ID
```http
GET /task/{id}
```

### Получение списка задач
```http
GET /task?status=pending
```

### Обновление статуса задачи
```http
PUT /task/{id}
Content-Type: application/json

{
  "status": "completed"
}
```

### Swagger документация
```http
GET /swagger/*
```

