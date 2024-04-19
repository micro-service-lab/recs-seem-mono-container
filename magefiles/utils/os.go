// Package utils This package defines utility functions for use in magefile.
package utils

import "runtime"

// IsWindows If the execution environment is Windows, return true.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
