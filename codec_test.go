package gsession_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opencomet-io/gsession"
)

func TestGobCodec(t *testing.T) {
	t.Parallel()

	codec := gsession.GobCodec{}

	t.Run("encoded data stays the same after decoding", func(t *testing.T) {
		t.Parallel()

		vals := map[string]any{
			"id":   25389,
			"role": "user",
		}
		expiry := time.Now().Add(1 * time.Hour).UTC()

		data, err := codec.Encode(vals, expiry)
		require.NoError(t, err)

		rVals, rExpiry, err := codec.Decode(data)
		require.NoError(t, err)

		assert.Equal(t, vals, rVals)
		assert.Equal(t, expiry, rExpiry)
	})

	t.Run("result of encoding differs based on the input values", func(t *testing.T) {
		t.Parallel()

		vals1 := map[string]any{"id": 12470}
		vals2 := map[string]any{"id": 791579}

		data1, err := codec.Encode(vals1, time.Time{})
		require.NoError(t, err)

		data2, err := codec.Encode(vals2, time.Time{})
		require.NoError(t, err)

		assert.NotEqual(t, data1, data2)
	})

	t.Run("result of encoding differs based on the input expiry", func(t *testing.T) {
		t.Parallel()

		ts1 := time.Now().Add(1 * time.Hour).UTC()
		ts2 := time.Now().Add(4 * time.Hour).UTC()

		data1, err := codec.Encode(map[string]any{}, ts1)
		require.NoError(t, err)

		data2, err := codec.Encode(map[string]any{}, ts2)
		require.NoError(t, err)

		assert.NotEqual(t, data1, data2)
	})
}
