package gsession

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomToken(t *testing.T) {
	t.Parallel()

	t.Run("generated token is different each time", func(t *testing.T) {
		t.Parallel()

		tokens := make(map[string]struct{})
		for i := 0; i < 8; i++ {
			token, err := generateRandomToken(32)
			require.NoError(t, err)

			_, ok := tokens[token]
			assert.Falsef(t, ok, "found two identical generated tokens: %s", token)

			tokens[token] = struct{}{}
		}
	})

	t.Run("token length too short", func(t *testing.T) {
		t.Parallel()

		_, err := generateRandomToken(0)
		assert.Error(t, err)

		_, err = generateRandomToken(-4)
		assert.Error(t, err)
	})

	tests := []struct {
		name            string
		tokenLength     int
		matchingPattern *regexp.Regexp
	}{
		{
			name:            "8 characters long random token",
			tokenLength:     8,
			matchingPattern: regexp.MustCompile(`^[a-zA-Z0-9]{8}$`),
		},
		{
			name:            "13 characters long random token",
			tokenLength:     13,
			matchingPattern: regexp.MustCompile(`^[a-zA-Z0-9]{13}$`),
		},
		{
			name:            "33 characters long random token",
			tokenLength:     33,
			matchingPattern: regexp.MustCompile(`^[a-zA-Z0-9]{33}$`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token, err := generateRandomToken(tt.tokenLength)
			require.NoError(t, err)

			match := tt.matchingPattern.MatchString(token)
			assert.True(t, match, "token string does not match the specified regexp pattern")
		})
	}
}
