package database

import (
	"crypto/rand"
	"fmt"
	"io"
)

// Thanks to https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb.
func GenerateRandomBytes(n int) ([]byte, error) {
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return nil, fmt.Errorf("crypto/rand is unavailable: Read() failed: %w", err)
	}

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// Note that err == nil only if we read len(b) bytes.
		return nil, fmt.Errorf("error while generating random bytes: %w", err)
	}

	return b, nil
}
