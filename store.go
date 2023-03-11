package gsession

import (
	"context"
	"time"
)

type Store interface {
	Get(ctx context.Context, token string) ([]byte, time.Time, error)
	Set(ctx context.Context, token string, data []byte, expiry time.Time) error
	Delete(ctx context.Context, token string) error
}
