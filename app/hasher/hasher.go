// Package hasher defines the hash interface and some implementations.
package hasher

// Hash ハッシュ。
type Hash interface {
	// Encrypt ハッシュ化する。
	Encrypt(text string) (string, error)
	// Compare ハッシュを比較する。
	Compare(text, hash string) error
}
