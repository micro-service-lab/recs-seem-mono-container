package service

import (
	"fmt"
	"math/rand"
)

// RandomColor generates a random color.
func RandomColor() string {
	return fmt.Sprintf("#%06X", rand.Intn(0xFFFFFF)) //nolint:gosec
}
