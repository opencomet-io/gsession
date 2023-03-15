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

func New() *InMemoryStore {
	s := &InMemoryStore{
		entries: make(map[string]entryPayload),
	}
	return s
}

func (s *InMemoryStore) Get(_ context.Context, token string) ([]byte, time.Time, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	payload, found := s.entries[token]
	if !found || time.Now().UTC().After(payload.expiry) {
		return nil, time.Time{}, false, nil
	}

	return payload.data, payload.expiry, true, nil
}

func (s *InMemoryStore) Set(_ context.Context, token string, data []byte, expiry time.Time) error {
	s.mu.Lock()

	payload := entryPayload{data, expiry}
	s.entries[token] = payload

	s.mu.Unlock()
	return nil
}

func (s *InMemoryStore) Delete(_ context.Context, token string) error {
	s.mu.Lock()

	delete(s.entries, token)

	s.mu.Unlock()
	return nil
}
