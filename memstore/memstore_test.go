package memstore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStore(t *testing.T) {
	t.Parallel()

	const (
		ExampleToken1 = "Token1"
		ExampleToken2 = "Token2"
		ExampleToken3 = "Token3"
	)

	var (
		ExampleData1 = []byte("Data 1")
		ExampleData2 = []byte("Data 2")
	)

	t.Run("adding new sessions to the store", func(t *testing.T) {
		t.Parallel()
		store := New()

		err := store.Set(context.Background(), ExampleToken1, ExampleData1, time.Now().Add(time.Hour).UTC())
		require.NoError(t, err)

		session, found := store.entries[ExampleToken1]
		require.True(t, found, "session not found")

		assert.False(t, session.expiry.IsZero(), "invalid session expiry time")
		assert.Equal(t, ExampleData1, session.data)
	})

	t.Run("session expiration", func(t *testing.T) {
		t.Parallel()
		store := New()

		err := store.Set(context.Background(), ExampleToken1, ExampleData1, time.Now().Add(-time.Hour).UTC())
		require.NoError(t, err)

		_, _, found, err := store.Get(context.Background(), ExampleToken1)
		require.NoError(t, err)

		assert.False(t, found, "retrieved expired session")
	})

	t.Run("deleting a single session", func(t *testing.T) {
		t.Parallel()
		store := New()
		store.entries[ExampleToken1] = entryPayload{ExampleData1, time.Now().Add(time.Hour).UTC()}
		store.entries[ExampleToken2] = entryPayload{ExampleData2, time.Now().Add(time.Hour).UTC()}

		err := store.Delete(context.Background(), ExampleToken2)
		require.NoError(t, err)

		_, found1 := store.entries[ExampleToken1]
		assert.True(t, found1, "cannot find a session entry that was not supposed to be deleted")

		_, found2 := store.entries[ExampleToken2]
		assert.False(t, found2, "found deleted session entry")
	})

	t.Run("deleting a non-existing session", func(t *testing.T) {
		t.Parallel()
		store := New()
		store.entries[ExampleToken1] = entryPayload{ExampleData1, time.Now().Add(time.Hour).UTC()}

		err := store.Delete(context.Background(), ExampleToken3)
		assert.NoError(t, err)

		_, found := store.entries[ExampleToken1]
		assert.True(t, found, "cannot find a session entry that was not supposed to be deleted")
	})

	t.Run("retrieving an existing session", func(t *testing.T) {
		t.Parallel()
		store := New()
		store.entries[ExampleToken1] = entryPayload{ExampleData1, time.Now().Add(time.Hour).UTC()}
		store.entries[ExampleToken2] = entryPayload{ExampleData2, time.Now().Add(time.Hour).UTC()}

		data, expiry, found, err := store.Get(context.Background(), ExampleToken1)
		require.NoError(t, err)
		require.True(t, found, "session not found")

		assert.False(t, expiry.IsZero(), "invalid session expiry time")
		assert.Equal(t, ExampleData1, data, "session data mismatch")
	})

	t.Run("retrieving a non-existing session", func(t *testing.T) {
		t.Parallel()
		store := New()
		store.entries[ExampleToken1] = entryPayload{ExampleData1, time.Now().Add(time.Hour).UTC()}

		_, _, found, err := store.Get(context.Background(), ExampleToken3)
		require.NoError(t, err)

		assert.False(t, found, "found a non-existing session")
	})

	t.Run("overriding session data", func(t *testing.T) {
		t.Parallel()
		store := New()
		store.entries[ExampleToken1] = entryPayload{ExampleData1, time.Now().Add(time.Hour).UTC()}

		want := []byte("Hello, world!")

		err := store.Set(context.Background(), ExampleToken1, want, time.Now().Add(time.Hour).UTC())
		require.NoError(t, err)

		payload, found := store.entries[ExampleToken1]
		require.True(t, found, "session not found")

		assert.Equal(t, want, payload.data, "session data was not overridden")
	})
}
