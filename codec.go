package gsession

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Codec interface {
	Encode(map[string]any, time.Time) ([]byte, error)
	Decode([]byte) (map[string]any, time.Time, error)
}

type carrier struct {
	Vals   map[string]any
	Expiry time.Time
}

type GobCodec struct{}

func (GobCodec) Encode(vals map[string]any, expiry time.Time) ([]byte, error) {
	c := carrier{vals, expiry}

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(&c); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (GobCodec) Decode(data []byte) (map[string]any, time.Time, error) {
	c := carrier{}

	r := bytes.NewReader(data)
	if err := gob.NewDecoder(r).Decode(&c); err != nil {
		return nil, time.Time{}, err
	}

	return c.Vals, c.Expiry, nil
}
