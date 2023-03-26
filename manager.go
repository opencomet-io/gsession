package gsession

import (
	"context"
	"time"
)

type Manager struct {
	Store       Store
	Codec       Codec
	TokenPrefix string
	TokenLength int
	Lifetime    time.Duration
}

func (m *Manager) insertSession(ctx context.Context, token string, vals map[string]any, expiry time.Time) error {
	storageToken := m.TokenPrefix + token

	data, err := m.Codec.Encode(vals, expiry)
	if err != nil {
		return err
	}

	return m.Store.Set(ctx, storageToken, data, expiry)
}

func (m *Manager) saveSession(ctx context.Context, token string, vals map[string]any) error {
	expiry := time.Now().Add(m.Lifetime).UTC()

	return m.insertSession(ctx, token, vals, expiry)
}

func (m *Manager) retrieveSession(ctx context.Context, token string) (map[string]any, time.Time, bool, error) {
	storageToken := m.TokenPrefix + token

	data, _, found, err := m.Store.Get(ctx, storageToken)
	if err != nil || !found {
		return nil, time.Time{}, false, err
	}

	vals, expiry, err := m.Codec.Decode(data)
	if err != nil {
		return nil, time.Time{}, false, err
	}

	return vals, expiry, true, nil
}

func (m *Manager) InitSession(ctx context.Context, vals map[string]any) (string, error) {
	token, err := generateRandomToken(m.TokenLength)
	if err != nil {
		return "", err
	}

	if vals == nil {
		vals = make(map[string]any)
	}

	err = m.saveSession(ctx, token, vals)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *Manager) GetSessionValues(ctx context.Context, token string) (map[string]any, error) {
	vals, _, found, err := m.retrieveSession(ctx, token)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNotFound
	}

	return vals, nil
}

func (m *Manager) SetSessionValues(ctx context.Context, token string, vals map[string]any) error {
	_, _, found, err := m.retrieveSession(ctx, token)
	if err != nil {
		return err
	}
	if !found {
		return ErrNotFound
	}

	return m.saveSession(ctx, token, vals)
}

func (m *Manager) AssertSessionValues(ctx context.Context, token string, want map[string]any) (bool, error) {
	return false, nil
}
