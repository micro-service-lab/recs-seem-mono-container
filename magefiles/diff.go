// Package main This is a task definition package for magefile, a task builder.
package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
	"github.com/micro-service-lab/recs-seem-mono-container/magefiles/utils"
)

// Diff is mage namespace for code diff code.
type Diff mg.Namespace

// Tabledoc diff table schema and table document.
func (s Diff) Tabledoc(ctx context.Context) {
	mg.CtxDeps(ctx, s.tabledoc)
}

// SQL diff sql query.
func (s Diff) SQL(ctx context.Context) {
	mg.CtxDeps(ctx, s.sql)
}

func (s Diff) tabledoc() error {
	repoRoot, err := utils.RepoRoot()
	if err != nil {
		return fmt.Errorf("get repo root: %w", err)
	}

	docDir := filepath.Join(repoRoot, "docs", "schema")

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("no .env file found: %w", err)
	}
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	env := map[string]string{
		"TBLS_DSN":      cfg.DBUrl,
		"TBLS_DOC_PATH": docDir,
	}

	if err := sh.RunWithV(env, "tbls", "diff"); err != nil {
		return fmt.Errorf("run diff table document: %w", err)
	}

	return nil
}

func (s Diff) sql() error {
	if err := sh.RunV("sqlc", "diff"); err != nil {
		return fmt.Errorf("run diff sql: %w", err)
	}

	return nil
}
