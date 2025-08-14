package task_repo

import (
	"context"
	"hash/fnv"
	"log/slog"
	"svc-task_master/src/domain"
	"sync"
	"time"
)

type SharderStorage struct {
	logger domain.ILogger
	Shard  []*Sharder
}

type Sharder struct {
	mu   sync.RWMutex
	Data map[string]*domain.Task
}

var _ domain.IInMemoRepository = &SharderStorage{}

func NewSharderStorage(numSharders int, ttl time.Duration, logger domain.ILogger) *SharderStorage {
	sharders := make([]*Sharder, numSharders)
	for i := 0; i < numSharders; i++ {
		sharders[i] = &Sharder{Data: make(map[string]*domain.Task)}
	}
	sharderStorage := &SharderStorage{
		logger: logger,
		Shard:  sharders,
	}
	if ttl > 0 {
		logger.Info("Starting TTL cleanup goroutine", slog.Attr{Key: "ttl", Value: slog.StringValue(ttl.String())})
		go sharderStorage.ClearForTTL(ttl)
	}
	return sharderStorage
}

func (s *SharderStorage) ClearForTTL(ttl time.Duration) {
	wg := sync.WaitGroup{}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		cutoff := now.Add(-ttl)
		deletedCount := 0

		for _, shard := range s.Shard {
			wg.Add(1)
			go func(sh *Sharder) {
				sh.mu.Lock()
				defer sh.mu.Unlock()
				defer wg.Done()

				for key, task := range sh.Data {
					if task.UpdatedAt.Before(cutoff) {
						delete(sh.Data, key)
						deletedCount++
					}
				}
			}(shard)
		}
		wg.Wait()

		if deletedCount > 0 {
			s.logger.Debug("Cleaned up expired tasks",
				slog.Attr{Key: "deleted_count", Value: slog.IntValue(deletedCount)},
				slog.Attr{Key: "cutoff_time", Value: slog.StringValue(cutoff.String())},
			)
		}
	}
}

func (s *SharderStorage) getSharder(key string) *Sharder {
	hashKey := fnv.New64a()
	hashKey.Write([]byte(key))
	sharderIndex := hashKey.Sum64() % uint64(len(s.Shard))
	return s.Shard[sharderIndex]
}

func (s *SharderStorage) GetAllFilterStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error) {
	s.logger.Debug("Getting all tasks with status filter",
		slog.Attr{Key: "status", Value: slog.StringValue(string(status))},
	)

	var result []domain.Task
	resultChan := make(chan []domain.Task, len(s.Shard))

	go func() {
		var wg sync.WaitGroup
		for _, shard := range s.Shard {
			wg.Add(1)
			go func(sh *Sharder) {
				defer wg.Done()

				sh.mu.RLock()
				defer sh.mu.RUnlock()

				var filtered []domain.Task
				for _, task := range sh.Data {
					select {
					case <-ctx.Done():
						s.logger.Warn("Context done while filtering tasks")
						return
					default:
						if status == "" || task.Status == status {
							filtered = append(filtered, *task)
						}
					}
				}
				resultChan <- filtered
			}(shard)
		}
		wg.Wait()
		close(resultChan)
	}()

	for {
		select {
		case <-ctx.Done():
			s.logger.Warn("Context done while collecting filtered tasks")
			return nil, ctx.Err()
		case filtered, ok := <-resultChan:
			if !ok {
				s.logger.Debug("Successfully retrieved filtered tasks",
					slog.Attr{Key: "count", Value: slog.IntValue(len(result))},
				)
				return result, nil
			}
			result = append(result, filtered...)
		}
	}
}

func (s *SharderStorage) Get(key string) (domain.Task, bool) {
	s.logger.Debug("Getting task by key",
		slog.Attr{Key: "key", Value: slog.StringValue(key)},
	)

	shard := s.getSharder(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	if data, ok := shard.Data[key]; ok {
		s.logger.Debug("Task found",
			slog.Attr{Key: "key", Value: slog.StringValue(key)},
			slog.Attr{Key: "status", Value: slog.StringValue(string(data.Status))},
		)
		return *data, ok
	}

	s.logger.Debug("Task not found",
		slog.Attr{Key: "key", Value: slog.StringValue(key)},
	)
	return domain.Task{}, false
}

func (s *SharderStorage) SetUpdate(key string, data domain.Task) {
	s.logger.Debug("Setting/updating task",
		slog.Attr{Key: "key", Value: slog.StringValue(key)},
		slog.Attr{Key: "status", Value: slog.StringValue(string(data.Status))},
	)

	shard := s.getSharder(key)
	shard.mu.Lock()
	shard.Data[key] = &data
	shard.mu.Unlock()
}

func (s *SharderStorage) UpdateStatus(key string, status domain.TaskStatus) {
	s.logger.Debug("Updating task status",
		slog.Attr{Key: "key", Value: slog.StringValue(key)},
		slog.Attr{Key: "new_status", Value: slog.StringValue(string(status))},
	)

	shard := s.getSharder(key)
	shard.mu.Lock()
	shard.Data[key].Status = status
	shard.Data[key].UpdatedAt = time.Now()
	shard.mu.Unlock()
}
