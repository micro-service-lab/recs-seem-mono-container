// Package app provides the application container.
package app

import (
	"context"
	"fmt"
	"os"

	"github.com/micro-service-lab/recs-seem-mono-container/app/hasher"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/pubsub"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store/pgadapter"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/ws"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock/fakeclock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// Container is a container for the application.
type Container struct {
	ServiceManager service.ManagerInterface
	Store          store.Store
	Storage        storage.Storage
	Clocker        clock.Clock
	Config         *config.Config
	Translator     i18n.Translation
	Hash           hasher.Hash
	SessionManager session.Manager
	Auth           auth.Auth
	Validator      validation.Validator
	PubsubService  pubsub.Service
	WebsocketHub   ws.HubInterface
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

	s3Bucket := os.Getenv("S3_BUCKET_NAME")
	if s3Bucket == "" {
		s3Bucket = "default-bucket"
	}
	s3CredentialKey := os.Getenv("S3_CREDENTIALS_KEY")
	if s3CredentialKey == "" {
		s3CredentialKey = "minio"
	}
	s3CredentialSecret := os.Getenv("S3_CREDENTIALS_SECRET")
	if s3CredentialSecret == "" {
		s3CredentialSecret = "minio123"
	}
	s3ExternalEndpoint := os.Getenv("S3_EXTERNAL_ENDPOINT")
	if s3ExternalEndpoint == "" {
		s3ExternalEndpoint = "http://localhost:9000"
	}
	s3Endpoint := fmt.Sprintf("%s:%d", cfg.StorageHost, cfg.StoragePort)

	s3, err := storage.NewS3(
		ctx,
		s3Endpoint,
		s3ExternalEndpoint,
		s3CredentialKey,
		s3CredentialSecret,
		s3Bucket,
	)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	c.Storage = s3

	h := hasher.NewBcrypt()

	c.Hash = h

	ssm, err := session.NewRedisManager(ctx, *cfg)
	if err != nil {
		return fmt.Errorf("failed to create session manager: %w", err)
	}

	c.SessionManager = ssm

	authSvc := auth.New(
		[]byte(cfg.AuthSecret),
		[]byte(cfg.AuthRefreshSecret),
		cfg.SecretIssuer,
		cfg.AuthAccessTokenExpiresIn, cfg.AuthRefreshTokenExpiresIn)
	vd, err := validation.NewRequestValidator()
	if err != nil {
		return fmt.Errorf("failed to create request validator: %w", err)
	}

	c.Auth = authSvc

	c.Validator = vd

	svc := service.NewManager(
		str, c.Translator, s3, h, clk, authSvc, ssm, *cfg,
	)

	c.ServiceManager = svc

	ps := pubsub.NewRedisService(*cfg)

	c.PubsubService = ps

	hub := ws.NewHub(ps)

	c.WebsocketHub = hub

	return nil
}

// Close closes the container.
func (c *Container) Close() error {
	if err := c.Store.Cleanup(context.Background()); err != nil {
		return fmt.Errorf("failed to cleanup store: %w", err)
	}

	return nil
}
