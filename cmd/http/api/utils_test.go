package api

import "testing"

func TestUUIDPath(t *testing.T) {
	//nolint: lll
	cases := map[string]struct {
		path     string
		expected string
	}{
		"simple": {
			path:     "/{abc:uuid}/bbb/{ccc:uuid}/ddd",
			expected: "/{abc:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/bbb/{ccc:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/ddd",
		},
		"last path is uuid": {
			path:     "/{abc:uuid}/bbb/{ccc:uuid}/{ddd:uuid}",
			expected: "/{abc:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/bbb/{ccc:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/{ddd:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}",
		},
		"no bracket": {
			path:     "/abc/def:uuid/ghi",
			expected: "/abc/def:uuid/ghi",
		},
		"plural depth bracket": {
			path:     "/{abc:{def:uuid}}",
			expected: "/{abc:{def:uuid}}",
		},
		"bracket in uuid": {
			path:     "/{abc:u{u}id}",
			expected: "/{abc:u{u}id}",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := uuidPath(c.path); got != c.expected {
				t.Errorf("unexpected result. expected: %s, but got: %s", c.expected, got)
			}
		})
	}
}
