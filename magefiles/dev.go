// Package main This is a task definition package for magefile, a task builder.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/micro-service-lab/recs-seem-mono-container/magefiles/utils"
)

const defaultDebuggingPort = 2345

// Dev starts development server with live reloading.
func Dev(ctx context.Context) {
	mg.CtxDeps(ctx, dev)
}

func dev() error {
	repoRoot, err := utils.RepoRoot()
	if err != nil {
		return fmt.Errorf("get repo root: %w", err)
	}

	port, err := debuggingPort()
	if err != nil {
		return fmt.Errorf("get debugging port: %w", err)
	}

	env := map[string]string{
		"GOFLAGS": "-buildvcs=false",
	}

	if err := sh.RunWithV(env, "arelo",
		"-t", repoRoot,
		"-p", "**/*.go",
		"-i", "**/*_test.go",
		"-i", "**/testutils/**",
		"-i", "**/testdata/**",
		"--",
		"dlv", "debug",
		"--accept-multiclient",
		"--api-version", "2",
		"--continue",
		"--headless",
		"--listen", fmt.Sprintf(":%d", port),
		"cmd/http/main.go",
		"cmd/http/server.go"); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

func debuggingPort() (uint16, error) {
	raw := os.Getenv("DEBUGGING_PORT")
	if raw == "" {
		return defaultDebuggingPort, nil
	}

	port, err := strconv.ParseUint(raw, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid debugging port: %w", err)
	}

	return uint16(port), nil
}
