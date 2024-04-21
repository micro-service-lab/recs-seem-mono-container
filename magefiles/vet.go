package main

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Vet is mage namespace for vet operations.
type Vet mg.Namespace

// SQL vet SQL.
func (s Vet) SQL(ctx context.Context) {
	mg.CtxDeps(ctx, s.sql)
}

func (s *Vet) sql(_ context.Context) error {
	if err := sh.RunV(
		"docker", "compose", "exec", "mono-api", "bash", "-c", "cd /app/server && sqlc vet",
	); err != nil {
		return fmt.Errorf("reset database: %w", err)
	}
	return nil
}
