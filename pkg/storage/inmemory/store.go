package inmemory

import (
	"context"
	"sync"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_store.go -package=${GOPACKAGE}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, bool)
	GetAll() map[string]interface{}
	Delete(ctx context.Context, key string) error
	IsEmpty() bool
}

type Store struct {
	mu     sync.RWMutex
	data   map[string]interface{}
	setOps map[string]chan struct{}
}

func New() Cache {
	return &Store{
		mu:   sync.RWMutex{},
		data: map[string]interface{}{},
	}
}

func (db *Store) GetAll() map[string]interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.data
}

func (db *Store) Get(_ context.Context, key string) (interface{}, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	item, ok := db.data[key]

	if !ok {
		return nil, false
	}

	return item, true
}

func (db *Store) Set(_ context.Context, key string, value interface{}) error {
	db.mu.Lock()
	db.data[key] = value
	db.mu.Unlock()

	return nil
}

func (db *Store) Delete(_ context.Context, key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)

	return nil
}

func (db *Store) IsEmpty() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return len(db.data) == 0
}
