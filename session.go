package gsession

import (
	"time"
)

type session struct {
	Token  string
	Values map[string]any
	Expiry time.Time
}
