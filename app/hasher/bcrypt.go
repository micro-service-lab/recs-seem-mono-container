package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt Bcryptハッシュ。
type Bcrypt struct{}

// NewBcrypt Bcryptハッシュを作成する。
func NewBcrypt() Bcrypt {
	return Bcrypt{}
}

// Encrypt ハッシュ化する。
func (b Bcrypt) Encrypt(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash from password: %w", err)
	}
	return string(hash), nil
}

// Compare ハッシュを比較する。
func (b Bcrypt) Compare(text, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text)); err != nil {
		return false, fmt.Errorf("failed to compare hash and text: %w", err)
	}
	return true, nil
}
