// Package factory provides a way to create test data for the tests.
package factory

import "sync"

// Factory is entity factory generator
type Factory struct {
	mu sync.Mutex
}

// Generator is a factory generator
var Generator Factory
