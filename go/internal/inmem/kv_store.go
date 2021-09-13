package inmem

import (
	"github.com/gc-plazas/kv-store/go/internal/errs"
	"sync"
)

type MemoryKeyValueStore struct {
	store map[string]string
	mu    sync.Mutex
}

func NewMemoryKeyValueStore() *MemoryKeyValueStore {
	return &MemoryKeyValueStore{store: make(map[string]string)}
}

func (m *MemoryKeyValueStore) Get(key string) (string, error) {
	if val, ok := m.store[key]; ok {
		return val, nil
	}
	return "", errs.ValueNotFound{}
}

func (m *MemoryKeyValueStore) Put(key, value string) error {
	m.mu.Lock()
	m.store[key] = value
	m.mu.Unlock()
	return nil
}
