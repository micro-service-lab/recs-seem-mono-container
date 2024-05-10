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

// Generate is mage namespace for code generation.
type Generate mg.Namespace

// Go run go generate.
func (s Generate) Go(ctx context.Context) {
	mg.CtxDeps(ctx, s.goGenerate)
}

// Tabledoc generates table document codes.
func (s Generate) Tabledoc(ctx context.Context) {
	mg.CtxDeps(ctx, s.tabledoc)
}

// Protoc generates go code for grpc.
func (s Generate) Protoc(ctx context.Context) {
	mg.CtxDeps(ctx, s.protoc)
}

// Migration generates migration file.
func (s Generate) Migration(
	ctx context.Context,
	service string,
) {
	mg.CtxDeps(ctx, s.migrationGenerator(service))
}

// SQL generates sql file.
func (s Generate) SQL(ctx context.Context) {
	mg.CtxDeps(ctx, s.sql)
}

func (s Generate) goGenerate() error {
	if err := sh.RunV("go", "generate", "./..."); err != nil {
		return fmt.Errorf("run go generate: %w", err)
	}

	return nil
}

func (s Generate) tabledoc() error {
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

	if err := sh.RunWithV(env, "tbls", "doc", "--rm-dist"); err != nil {
		return fmt.Errorf("run generate table document: %w", err)
	}

	return nil
}

func (s Generate) protoc() error {
	repoRoot, err := utils.RepoRoot()
	if err != nil {
		return fmt.Errorf("get repo root: %w", err)
	}

	protoDir := filepath.Join(repoRoot, "proto")

	genCmd := fmt.Sprintf("protoc -I%[1]s --go_out=%[1]s --go-grpc_out=%[1]s %[2]s/*.proto", repoRoot, protoDir)

	if err := sh.RunV("bash",
		"-c",
		genCmd,
	); err != nil {
		return fmt.Errorf("run generate go code for grpc: %w", err)
	}

	return nil
}

func (s Generate) migrationGenerator(
	filename string,
) func() error {
	return func() error {
		repoRoot, err := utils.RepoRoot()
		if err != nil {
			return fmt.Errorf("get repo root: %w", err)
		}

		migrateDir := filepath.Join(repoRoot, "db", "migrations")

		if err := sh.RunV("migrate",
			"create", "-ext", "sql", "-dir", migrateDir, "-seq", filename,
		); err != nil {
			return fmt.Errorf("run generate migration: %w", err)
		}

		return nil
	}
}

func (s Generate) sql() error {
	genCmd := "sqlc generate"

	if err := sh.RunV("bash",
		"-c",
		genCmd,
	); err != nil {
		return fmt.Errorf("run generate go code for sqlc: %w", err)
	}

	return nil
}
