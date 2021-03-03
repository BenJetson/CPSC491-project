package app

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// A Password is a wrapper type for bcrypt password hashes.
type Password string

// NewPassword creates a new hashed Password given its plaintext.
func NewPassword(plaintext string) (Password, error) {
	plainBytes := []byte(plaintext)
	hashBytes, err := bcrypt.GenerateFromPassword(
		plainBytes,
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	return Password(hashBytes), nil
}

// Verify verifies a hashed Password against a given plaintext.
func (p Password) Verify(plaintext string) bool {
	hashBytes := []byte(p)
	plainBytes := []byte(plaintext)

	// CompareHashAndPassword returns nil error when the plaintext and hash
	// match, and an error otherwise.
	return bcrypt.CompareHashAndPassword(hashBytes, plainBytes) == nil
}
