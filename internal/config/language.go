package config

import (
	"errors"
)

// Language Indicates the language setting that can be specified.
type Language string

const (
	// Japanese ja_JP language
	Japanese Language = "ja"
	// English en language
	English Language = "en"
)

// parseLanguage Parses the language setting string.
func parseLanguage(v string) (any, error) {
	if v == "" {
		return Language(""), errors.New("invalid value")
	}

	switch lng := Language(v); lng {
	case Japanese, English:
		return lng, nil
	default:
		return nil, errors.New("invalid value")
	}
}
