package gsession

import (
	"time"
)

type kindOfSession string

const (
	tempSession    kindOfSession = "temp"
	refreshSession kindOfSession = "refresh"
)

type session struct {
	Token  string
	Kind   kindOfSession
	Values map[string]any
	Expiry time.Time
}
