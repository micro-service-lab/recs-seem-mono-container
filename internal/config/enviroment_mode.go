package config

import (
	"fmt"
)

// EnvironmentMode Indicates the setting of the environment mode that can be specified.
type EnvironmentMode string

const (
	// ProductionEnv Indicates production environment mode.
	ProductionEnv EnvironmentMode = "production"
	// StagingEnv Indicates staging environment mode.
	StagingEnv EnvironmentMode = "staging"
	// TestEnv Indicates the test environment mode.
	TestEnv EnvironmentMode = "test"
	// DevelopmentEnv Indicates development (local) environment mode.
	DevelopmentEnv EnvironmentMode = "development"
)

// parseEnvironmentMode Parses the environment mode setting string.
func parseEnvironmentMode(v string) (any, error) {
	if v == "" {
		return ProductionEnv, nil
	}

	switch ev := EnvironmentMode(v); ev {
	case ProductionEnv, StagingEnv, TestEnv, DevelopmentEnv:
		return ev, nil
	default:
		return nil, fmt.Errorf("specify one of the following items: %s, %s, %s, %s",
			ProductionEnv, StagingEnv, TestEnv, DevelopmentEnv)
	}
}
