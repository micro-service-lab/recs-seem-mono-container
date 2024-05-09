// Package batch contains the batch processing logic for the application.
package batch

import (
	"context"
	"reflect"

	"github.com/google/uuid"
)

// Batch is a interface for batch.
type Batch interface {
	// Run runs the batch.
	Run(ctx context.Context) error
	// RunDiff runs the batch only if there is a difference.
	RunDiff(ctx context.Context, notDel, deepEqual bool) error
}

// return match index
// if not found, return -1
func contains(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func isDeepEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func removeUUID(s []uuid.UUID, i int) ([]uuid.UUID, uuid.UUID) {
	if len(s) < i {
		return s, uuid.Nil
	}
	u := s[i]
	s[i] = s[len(s)-1]
	return s[:len(s)-1], u
}

func removeString(s []string, i int) ([]string, string) { //nolint:unparam
	if len(s) < i {
		return s, ""
	}
	u := s[i]
	s[i] = s[len(s)-1]
	return s[:len(s)-1], u
}
