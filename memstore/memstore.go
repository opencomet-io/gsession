package memstore

import (
	"context"
	"sync"
	"time"

	"github.com/opencomet-io/gsession"
)

type entryPayload struct {
	data   []byte
	expiry time.Time
}

type InMemoryStore struct {
	entries map[string]entryPayload
	mu      sync.RWMutex
}

var _ gsession.Store = (*InMemoryStore)(nil)

func (s *InMemoryStore) Get(ctx context.Context, token string) ([]byte, time.Time, error) {
	return nil, time.Time{}, nil
}

func (s *InMemoryStore) Set(ctx context.Context, token string, data []byte, expiry time.Time) error {
	return nil
}

func (s *InMemoryStore) Delete(ctx context.Context, token string) error {
	return nil
}
