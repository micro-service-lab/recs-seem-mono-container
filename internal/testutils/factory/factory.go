// Package factory provides a way to create test data for the tests.
package factory

// DataFactory is an interface for generating test data.
type DataFactory interface {
	// Generate generates test data.
	Generate() error
}
