package task_repo

import (
	"context"
	"hash/fnv"
	"svc-task_master/src/domain"
	"sync"
	"time"
)

type SharderStorage struct {
	Shard []*Sharder
}

type Sharder struct {
	mu   sync.RWMutex
	Data map[string]*domain.Task
}

var _ domain.IInMemoRepository = &SharderStorage{}

func NewSharderStorage(numSharders int) *SharderStorage {
	sharders := make([]*Sharder, numSharders)
	for i := 0; i < numSharders; i++ {
		sharders[i] = &Sharder{Data: make(map[string]*domain.Task)}
	}
	return &SharderStorage{Shard: sharders}
}

func (s *SharderStorage) getSharder(key string) *Sharder {
	hashKey := fnv.New64a()
	hashKey.Write([]byte(key))
	sharderIndex := hashKey.Sum64() % uint64(len(s.Shard))
	return s.Shard[sharderIndex]
}

func (s *SharderStorage) GetAllFilterStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error) {
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
			return nil, ctx.Err()
		case filtered, ok := <-resultChan:
			if !ok {
				return result, nil
			}
			result = append(result, filtered...)
		}
	}
}

func (s *SharderStorage) Get(key string) (domain.Task, bool) {
	shard := s.getSharder(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	if data, ok := shard.Data[key]; ok {
		return *data, ok
	}
	return domain.Task{}, false
}

func (s *SharderStorage) SetUpdate(key string, data domain.Task) {
	shard := s.getSharder(key)
	shard.mu.Lock()
	shard.Data[key] = &data
	shard.mu.Unlock()
}

func (s *SharderStorage) UpdateStatus(key string, status domain.TaskStatus) {
	shard := s.getSharder(key)
	shard.mu.Lock()
	shard.Data[key].Status = status
	shard.Data[key].UpdatedAt = time.Now()
	shard.mu.Unlock()
}
