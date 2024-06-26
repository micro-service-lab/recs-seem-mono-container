package main

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/micro-service-lab/recs-seem-mono-container/magefiles/utils"
)

// App is mage namespace for app operations.
type App mg.Namespace

// Dev starts development server with live reloading.
func (a App) Dev(ctx context.Context) {
	mg.CtxDeps(ctx, a.dev)
}

// Kill development server.
func (a App) Kill(ctx context.Context) {
	mg.CtxDeps(ctx, a.kill)
}

// Serve starts the application.
func (a App) Serve(ctx context.Context) {
	mg.CtxDeps(ctx, a.serve)
}

// Up starts the application.
func (a App) Up(ctx context.Context) {
	mg.CtxDeps(ctx, a.up)
}

// Down stops the application.
func (a App) Down(ctx context.Context) {
	mg.CtxDeps(ctx, a.down)
}

// Ps shows the status of the application.
func (a App) Ps(ctx context.Context) {
	mg.CtxDeps(ctx, a.ps)
}

// Migrate migrates the database schema.
func (a App) Migrate(ctx context.Context) {
	mg.CtxDeps(ctx, a.migrate)
}

// Rollback rolls back the database schema.
func (a App) Rollback(ctx context.Context) {
	mg.CtxDeps(ctx, a.rollback)
}

// Seed loads seed data.(options can be specified by colon-separated)
func (a App) Seed(ctx context.Context, target string) {
	mg.CtxDeps(ctx, a.seedGenerator(target))
}

// Tabledoc generates table documentation.
func (a App) Tabledoc(ctx context.Context) {
	mg.CtxDeps(ctx, a.tabledoc)
}

// Force forces version.
func (a App) Force(ctx context.Context, version string) {
	mg.CtxDeps(ctx, a.forceVersionGenerator(version))
}

// Create creates the database.
func (a App) Create(ctx context.Context) {
	mg.CtxDeps(ctx, a.create)
}

// Drop deletes the database.
func (a App) Drop(ctx context.Context) {
	mg.CtxDeps(ctx, a.drop)
}

// Reset resets the database.
func (a App) Reset(ctx context.Context) {
	mg.CtxDeps(ctx, a.reset)
}

// Log shows the logs of the service.
func (a App) Log(ctx context.Context, service string) {
	mg.CtxDeps(ctx, a.logGenerator(service))
}

// Bash opens a bash session in the service container.
func (a App) Bash(ctx context.Context, service string) {
	mg.CtxDeps(ctx, a.bashGenerator(service))
}

// Clean removes all __debug_bin files.
func (a App) Clean(ctx context.Context) {
	mg.CtxDeps(ctx, a.clean)
}

func (a App) dev() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "dev"); err != nil {
		return fmt.Errorf("run server: %w", err)
	}
	return nil
}

func (a App) kill() error {
	if err := sh.RunV("docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "kill"); err != nil {
		return fmt.Errorf("kill container: %w", err)
	}
	return nil
}

func (a App) serve() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "serve"); err != nil {
		return fmt.Errorf("run server: %w", err)
	}
	return nil
}

func (a App) up() error {
	if err := sh.RunV("docker", "compose", "up", "-d"); err != nil {
		return fmt.Errorf("up container: %w", err)
	}

	return nil
}

func (a App) down() error {
	if err := sh.RunV("docker", "compose", "down"); err != nil {
		return fmt.Errorf("down container: %w", err)
	}
	return nil
}

func (a App) ps() error {
	if err := sh.RunV("docker", "compose", "ps"); err != nil {
		return fmt.Errorf("ps container: %w", err)
	}
	return nil
}

func (a App) create() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "db:create"); err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	return nil
}

func (a App) migrate() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "db:migrate"); err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	return nil
}

func (a App) rollback() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "db:rollback"); err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	return nil
}

func (a App) seedGenerator(target string) func() error {
	return func() error {
		if err := sh.RunV(
			"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "db:seed", target); err != nil {
			return fmt.Errorf("create database: %w", err)
		}
		return nil
	}
}

func (a App) tabledoc() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "generate:tabledoc"); err != nil {
		return fmt.Errorf("generate table documentation: %w", err)
	}

	return nil
}

func (a App) forceVersionGenerator(version string) func() error {
	return func() error {
		if err := sh.RunV(
			"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "db:force", version); err != nil {
			return fmt.Errorf("force version: %w", err)
		}
		return nil
	}
}

func (a App) drop() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "db:drop"); err != nil {
		return fmt.Errorf("drop database: %w", err)
	}
	return nil
}

func (a App) reset() error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "mage", "-d", "/app/server", "db:reset"); err != nil {
		return fmt.Errorf("reset database: %w", err)
	}
	return nil
}

func (a App) logGenerator(service string) func() error {
	return func() error {
		if err := sh.RunV("docker", "compose", "logs", "--tail", "400", "-f", service); err != nil {
			return fmt.Errorf("log container: %w", err)
		}
		return nil
	}
}

func (a App) bashGenerator(service string) func() error {
	return func() error {
		if err := sh.RunV("docker", "compose", "exec", service, "bash"); err != nil {
			return fmt.Errorf("bash container: %w", err)
		}
		return nil
	}
}

func (a App) clean() error {
	rootDir, err := utils.RepoRoot()
	if err != nil {
		return fmt.Errorf("get repo root: %w", err)
	}
	err = utils.RemoveDebugBinFiles(rootDir)
	if err != nil {
		return fmt.Errorf("remove debug bin files: %w", err)
	}
	return nil
}
