// Package main This is a task definition package for magefile, a task builder.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/micro-service-lab/recs-seem-mono-container/magefiles/utils"
)

// Test is mage namespace for code test.
type Test mg.Namespace

// All execute all test.
func (s Test) All(ctx context.Context) {
	mg.CtxDeps(ctx, s.all)
}

// Simple execute simple test.
func (s Test) Simple(ctx context.Context) {
	mg.CtxDeps(ctx, s.simple)
}

// Short execute short test.
func (s Test) Short(ctx context.Context) {
	mg.CtxDeps(ctx, s.short)
}

// Race execute with checking concurrent access.
func (s Test) Race(ctx context.Context) {
	mg.CtxDeps(ctx, s.race)
}

// Cover execute with coverage.
func (s Test) Cover(ctx context.Context) {
	mg.CtxDeps(ctx, s.cover)
}

func (s Test) all() error {
	testTarget, err := getAllTestTarget()
	if err != nil {
		return fmt.Errorf("get all test target: %w", err)
	}
	if ran, err := sh.Exec(
		map[string]string{},
		os.Stdout,
		os.Stderr,
		"go",
		"test",
		"-cover",
		"-race",
		"-v",
		testTarget,
	); !ran || err != nil {
		return fmt.Errorf("run test: %w", err)
	}

	return nil
}

func (s Test) simple() error {
	testTarget, err := getAllTestTarget()
	if err != nil {
		return fmt.Errorf("get all test target: %w", err)
	}
	if ran, err := sh.Exec(
		map[string]string{},
		os.Stdout,
		os.Stderr,
		"go",
		"test",
		testTarget,
	); !ran || err != nil {
		return fmt.Errorf("run test: %w", err)
	}

	return nil
}

func (s Test) short() error {
	testTarget, err := getAllTestTarget()
	if err != nil {
		return fmt.Errorf("get all test target: %w", err)
	}
	if ran, err := sh.Exec(
		map[string]string{},
		os.Stdout,
		os.Stderr,
		"go",
		"test",
		"-short",
		testTarget,
	); !ran || err != nil {
		return fmt.Errorf("run test: %w", err)
	}

	return nil
}

func (s Test) race() error {
	testTarget, err := getAllTestTarget()
	if err != nil {
		return fmt.Errorf("get all test target: %w", err)
	}
	if ran, err := sh.Exec(
		map[string]string{},
		os.Stdout,
		os.Stderr,
		"go",
		"test",
		"-race",
		testTarget,
	); !ran || err != nil {
		return fmt.Errorf("run test: %w", err)
	}

	return nil
}

func (s Test) cover() error {
	testTarget, err := getAllTestTarget()
	if err != nil {
		return fmt.Errorf("get all test target: %w", err)
	}
	if ran, err := sh.Exec(
		map[string]string{},
		os.Stdout,
		os.Stderr,
		"go",
		"test",
		"-cover",
		testTarget,
	); !ran || err != nil {
		return fmt.Errorf("run test: %w", err)
	}

	return nil
}

func getAllTestTarget() (string, error) {
	repoRoot, err := utils.RepoRoot()
	if err != nil {
		return "", fmt.Errorf("get repo root: %w", err)
	}
	testTarget := fmt.Sprintf("%s/...", repoRoot)
	return testTarget, nil
}
