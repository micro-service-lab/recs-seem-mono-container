// Package app provides the application container.
package app

import (
	"context"
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store/pgadapter"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
)

// Container is a container for the application.
type Container struct {
	ServiceManager service.ManagerInterface
	Store          store.Store
	Clocker        clock.Clock
	Config         *config.Config
	Translator     i18n.Translation
}

// NewContainer creates a new Container.
func NewContainer() *Container {
	return &Container{}
}

// Init resolves dependencies.
func (c *Container) Init(ctx context.Context) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	c.Config = cfg

	clk := clock.New()
	if cfg.FakeTime.Enabled {
		fakeClk := fakeclock.New(cfg.FakeTime.Time)
		clk = fakeClk
	}

	c.Clocker = clk

	str, err := pgadapter.NewPgAdapter(
		ctx,
		clk,
		*cfg,
	)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	c.Store = str

	c.Translator, err = i18n.NewTranslator()
	if err != nil {
		return fmt.Errorf("failed to create translator: %w", err)
	}

	svc := service.NewManager(str, c.Translator)

	c.ServiceManager = svc

	return nil
}

// Close closes the container.
func (c *Container) Close() error {
	if err := c.Store.Cleanup(context.Background()); err != nil {
		return fmt.Errorf("failed to cleanup store: %w", err)
	}

	return nil
}
