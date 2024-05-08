// Package batch contains the batch processing logic for the application.
package batch

import (
	"reflect"

	"github.com/google/uuid"
)

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
