package main

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/database"
	"github.com/micro-service-lab/recs-seem-mono-container/magefiles/utils"
)

// DB is mage namespace for db operations.
type DB mg.Namespace

// Create creates database.
func (s DB) Create(ctx context.Context) {
	mg.CtxDeps(ctx, s.create)
}

// Migrate migrates database schema.
func (s DB) Migrate(ctx context.Context) {
	mg.CtxDeps(ctx, s.migrate)
}

// Rollback rolls back database schema.
func (s DB) Rollback(ctx context.Context) {
	mg.CtxDeps(ctx, s.rollback)
}

// Seed loads seed data.(options can be specified by colon-separated)
func (s DB) Seed(ctx context.Context, target string) {
	targetsArr := strings.Split(target, ":")
	mg.CtxDeps(ctx, s.seedGenerator(targetsArr...))
}

// Force forces version.
func (s DB) Force(ctx context.Context, version string) {
	mg.CtxDeps(ctx, s.forceVersionGenerator(version))
}

// Drop deletes database.
func (s DB) Drop(ctx context.Context) {
	mg.CtxDeps(ctx, s.drop)
}

// Reset resets database.
func (s DB) Reset(ctx context.Context) {
	mg.SerialCtxDeps(ctx, s.drop, s.create, s.migrate, s.seedGenerator("all", "-f"))
}

// Fake inserts fake data.
func (s DB) Fake(ctx context.Context) {
	mg.CtxDeps(ctx, s.fake)
}

func (s *DB) create(ctx context.Context) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	db, err := database.Open(cfg.DBHost, uint16(cfg.DBPort), cfg.DBName, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close(ctx)

	query := fmt.Sprintf("CREATE DATABASE `%s`", cfg.DBName)

	if _, err := db.Exec(ctx, query); err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	fmt.Printf("database %q has been created.\n", cfg.DBName)

	return nil
}

func (s *DB) migrate() error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	repoRoot, err := utils.RepoRoot()
	if err != nil {
		return fmt.Errorf("get repo root: %w", err)
	}

	dbDir := filepath.Join(repoRoot, "db")

	args := []string{
		"--path", filepath.Join(dbDir, "migrations"),
		"--database", cfg.DBUrl,
		"-verbose",
		"up",
	}
	if err := sh.RunV("migrate", args...); err != nil {
		return fmt.Errorf("run migrate: %w", err)
	}

	return nil
}

func (s *DB) rollback() error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	repoRoot, err := utils.RepoRoot()
	if err != nil {
		return fmt.Errorf("get repo root: %w", err)
	}

	dbDir := filepath.Join(repoRoot, "db")

	args := []string{
		"--path", filepath.Join(dbDir, "migrations"),
		"--database", cfg.DBUrl,
		"-verbose",
		"down",
	}
	if err := sh.RunV("migrate", args...); err != nil {
		return fmt.Errorf("run migrate: %w", err)
	}
	return nil
}

func (s *DB) seedGenerator(targets ...string) func() error {
	return func() error {
		repoRoot, err := utils.RepoRoot()
		if err != nil {
			return fmt.Errorf("get repo root: %w", err)
		}

		cmdFile := filepath.Join(repoRoot, "cmd", "cli", "main.go")

		args := []string{
			"run", cmdFile, "seed",
		}
		args = append(args, targets...)

		if err := sh.RunV("go", args...); err != nil {
			return fmt.Errorf("execute seed file: %w", err)
		}

		return nil
	}
}

func (s *DB) forceVersionGenerator(version string) func() error {
	return func() error {
		cfg, err := config.Get()
		if err != nil {
			return fmt.Errorf("get config: %w", err)
		}

		repoRoot, err := utils.RepoRoot()
		if err != nil {
			return fmt.Errorf("get repo root: %w", err)
		}

		dbDir := filepath.Join(repoRoot, "db")

		args := []string{
			"--path", filepath.Join(dbDir, "migrations"),
			"--database", cfg.DBUrl,
			"force", version,
		}
		if err := sh.RunV("migrate", args...); err != nil {
			return fmt.Errorf("run migrate: %w", err)
		}

		return nil
	}
}

func (s *DB) drop(ctx context.Context) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	db, err := database.Open(cfg.DBHost, uint16(cfg.DBPort), cfg.DBName, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close(ctx)

	query := fmt.Sprintf("DROP DATABASE %s", cfg.DBName)
	if _, err := db.Exec(ctx, query); err != nil {
		return fmt.Errorf("drop database: %w", err)
	}

	fmt.Printf("database %q has been deleted.\n", cfg.DBName)

	return nil
}

const fakedataNum = 300

func (s *DB) fake(ctx context.Context) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	db, err := database.Open(cfg.DBHost, uint16(cfg.DBPort), cfg.DBName, cfg.DBUsername, cfg.DBPassword)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close(ctx)

	users := newFakeUsers(fakedataNum)

	var query strings.Builder
	//nolint:lll
	query.WriteString("INSERT INTO users (name, play_count, high_score, max_depth, last_logged_in_at, last_played_at, created_at, updated_at) ")
	args := make([]any, 0, fakedataNum*9) //nolint:gomnd
	for i, user := range users {
		if i != 0 {
			query.WriteString(" UNION ")
		}
		query.WriteString("SELECT ?, ?, ?, ?, ?, ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM users WHERE name = ?)")

		args = append(args,
			user.name, user.playCount, user.highScore, user.maxDepth,
			user.lastLoggedInAt, user.lastPlayedAt, user.createdAt, user.updatedAt,
			user.name,
		)
	}

	if _, err := db.Exec(ctx, query.String(), args...); err != nil {
		return fmt.Errorf("insert fake data: %w", err)
	}

	fmt.Println("fake data have been created.")

	return nil
}

const (
	minCount    = 1
	maxCount    = 100
	minScore    = 0
	maxScore    = 10000
	minDepth    = 0
	maxDepth    = 2000
	maxInterval = 3600
)

type fakeUser struct {
	name           string
	playCount      int32
	highScore      int32
	maxDepth       int32
	lastLoggedInAt time.Time
	lastPlayedAt   time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

//nolint:gosec
func newFakeUsers(num int) []*fakeUser {
	userMap := make(map[string]*fakeUser, num)
	for i := 0; i < num; i++ {
		var name string
		for {
			name = genFakeUsername()
			if _, ok := userMap[name]; !ok {
				break
			}
		}

		count := rand.Int31n(maxCount + 1)
		if count < minCount {
			count = minCount
		}

		score := rand.Int31n(maxScore + 1)
		if score < minScore {
			score = minScore
		}

		depth := rand.Int31n(maxDepth + 1)
		if depth < minDepth {
			depth = minDepth
		}

		updated := time.Now().Add(-time.Duration(rand.Int31n(maxInterval)) * time.Second)
		played := updated
		loggedIn := played.Add(-time.Duration(rand.Int31n(maxInterval)) * time.Second)
		created := loggedIn.Add(-time.Duration(rand.Int31n(maxInterval)) * time.Second)

		userMap[name] = &fakeUser{
			name:           name,
			playCount:      count,
			highScore:      score,
			maxDepth:       depth,
			lastLoggedInAt: loggedIn,
			lastPlayedAt:   played,
			createdAt:      created,
			updatedAt:      updated,
		}
	}

	users := make([]*fakeUser, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, user)
	}

	return users
}

var fakeUserNames = []string{
	"account",
	"com",
	"computer",
	"data",
	"dummy",
	"example",
	"fake",
	"fakedata",
	"fakeuser",
	"npc",
	"player",
	"record",
	"test",
	"testdata",
	"testuser",
	"unknown",
	"user",
}

var suffixLetters = "abcdefghijklmnopqrstuvwxyz0123456789"

//nolint:gosec
func genFakeUsername() string {
	base := fakeUserNames[rand.Intn(len(fakeUserNames))]
	num := rand.Intn(10) //nolint:gomnd
	suffix := suffixLetters[rand.Intn(len(suffixLetters))]
	return fmt.Sprintf("%s%d%s", base, num, string(suffix))
}
