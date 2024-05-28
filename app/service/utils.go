package service

import (
	"fmt"
	"math/rand"
)

// RandomColor generates a random color.
func RandomColor() string {
	return fmt.Sprintf("#%06X", rand.Intn(0xFFFFFF)) //nolint:gosec,gomnd
}

// ArrayShuffle shuffles an array.
func ArrayShuffle[T any](arr []T) {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}
