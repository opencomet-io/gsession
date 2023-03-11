package gsession

import (
	"time"
)

type kindOfSession string

const (
	TempSession    kindOfSession = "temp"
	RefreshSession kindOfSession = "refresh"
)

type session struct {
	Token  string
	Kind   kindOfSession
	Values map[string]any
	Expiry time.Time
}
