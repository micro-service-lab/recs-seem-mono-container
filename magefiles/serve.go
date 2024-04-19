package main

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Serve starts development server.
func Serve(ctx context.Context) {
	mg.CtxDeps(ctx, serve)
}

func serve() error {
	env := map[string]string{}

	if err := sh.RunWithV(env, "go", "run", "."); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}
