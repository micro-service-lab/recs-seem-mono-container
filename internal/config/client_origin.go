package config

import "strings"

// ClientOrigin indicates the setting of the client origin.
type ClientOrigin []string

// parseClientOrigin Parses the client origin setting string.
func parseClientOrigin(v string) (any, error) {
	if v == "" {
		return nil, nil
	}
	// split the string by comma
	cors := strings.Split(v, ",")
	// remove empty strings
	for i := 0; i < len(cors); i++ {
		cors[i] = strings.TrimSpace(cors[i])
		if cors[i] == "" {
			cors = append(cors[:i], cors[i+1:]...)
			i--
		}
	}
	// return the split string
	return ClientOrigin(cors), nil
}
