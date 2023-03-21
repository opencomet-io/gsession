package gsession

import (
	"time"
)

type session struct {
	Token     string
	IsRefresh bool
	Values    map[string]any
	Expiry    time.Time
}
