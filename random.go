package gsession

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

const TokenCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomToken(length int) (string, error) {
	if length < 1 {
		return "", errors.New("token length should be greater than 0")
	}

	t := make([]uint8, length)
	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(TokenCharset))))
		if err != nil {
			return "", fmt.Errorf("an error occurred while generating a random token: %w", err)
		}
		t[i] = TokenCharset[charIndex.Int64()]
	}
	return string(t), nil
}
