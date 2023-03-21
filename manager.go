package gsession

import (
	"context"
	"time"
)

type RefreshOpts struct {
	Store       Store
	TokenPrefix string
	Lifetime    time.Duration
}

type Manager struct {
	Store       Store
	Codec       Codec
	TokenLength int
	Lifetime    time.Duration
	Refresh     RefreshOpts
}

func (m *Manager) saveSession(ctx context.Context, token string, vals map[string]any) error {
	expiry := time.Now().Add(m.Lifetime).UTC()

	data, err := m.Codec.Encode(vals, expiry)
	if err != nil {
		return err
	}

	return m.Store.Set(ctx, token, data, expiry)
}

func (m *Manager) saveRefreshSession(ctx context.Context, token string, vals map[string]any, expiry time.Time) error {
	storageToken := m.Refresh.TokenPrefix + token

	data, err := m.Codec.Encode(vals, expiry)
	if err != nil {
		return err
	}

	return m.Refresh.Store.Set(ctx, storageToken, data, expiry)
}

func (m *Manager) InitSession(ctx context.Context) (string, error) {
	token, err := generateRandomToken(m.TokenLength)
	if err != nil {
		return "", err
	}

	vals := make(map[string]any)
	err = m.saveSession(ctx, token, vals)
	if err != nil {
		return "", err
	}

	return token, nil
}
