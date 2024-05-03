package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
)

func TestGet(t *testing.T) {
	// Unset environment variables for test
	envKeys := []string{
		"PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_URL",
		"AUTH_SECRET",
		"SECRET_ISSUER",
		"CLIENT_ORIGIN",
		"DEBUG_CORS",
		"APP_DEBUG",
		"APP_ENV",
		"FAKE_TIME",
		"LOG_LEVEL",
	}
	for _, v := range envKeys {
		t.Setenv(v, "")
		os.Unsetenv(v)
	}

	cases := []struct {
		name   string
		env    map[string]string
		out    *config.Config
		failed bool
	}{
		{
			name: "minimum",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AppEnv:       config.ProductionEnv,
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				LogLevel:     config.InfoLevel,
			},
		},
		{
			name: "full",
			env: map[string]string{
				"PORT":          "3000",
				"DB_HOST":       "db",
				"DB_PORT":       "9999",
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_PASSWORD":   "password",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"CLIENT_ORIGIN": "http://localhost:8080,http://localhost:8081",
				"DEBUG_CORS":    "true",
				"APP_DEBUG":     "true",
				"APP_ENV":       "development",
				"FAKE_TIME":     "true",
				"LOG_LEVEL":     "debug",
			},
			out: &config.Config{
				Port:         3000,
				DBHost:       "db",
				DBPort:       9999,
				DBName:       "app",
				DBUsername:   "user",
				DBPassword:   "password",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				ClientOrigin: config.ClientOrigin{"http://localhost:8080", "http://localhost:8081"},
				DebugCORS:    true,
				AppDebug:     true,
				AppEnv:       config.DevelopmentEnv,
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    config.DefaultFakeTime,
				},
				LogLevel: config.DebugLevel,
			},
		},
		{
			name: "FAKE_TIME is RFC3339 string",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"FAKE_TIME":     "2023-01-02T12:34:56Z",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				AppEnv:       config.ProductionEnv,
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    time.Date(2023, 1, 2, 12, 34, 56, 0, time.UTC),
				},
				LogLevel: config.InfoLevel,
			},
		},
		{
			name: "FAKE_TIME is true",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"FAKE_TIME":     "true",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				AppEnv:       config.ProductionEnv,
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    config.DefaultFakeTime,
				},
				LogLevel: config.InfoLevel,
			},
		},
		{
			name: "FAKE_TIME is 1",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"FAKE_TIME":     "1",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				AppEnv:       config.ProductionEnv,
				FakeTime: config.FakeTimeMode{
					Enabled: true,
					Time:    config.DefaultFakeTime,
				},
				LogLevel: config.InfoLevel,
			},
		},
		{
			name: "FAKE_TIME is false",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"FAKE_TIME":     "false",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				AppEnv:       config.ProductionEnv,
				LogLevel:     config.InfoLevel,
			},
		},
		{
			name: "FAKE_TIME is 0",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"FAKE_TIME":     "0",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				AppEnv:       config.ProductionEnv,
				LogLevel:     config.InfoLevel,
			},
		},
		{
			name: "FAKE_TIME is empty string",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"FAKE_TIME":     "",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				AppEnv:       config.ProductionEnv,
				LogLevel:     config.InfoLevel,
			},
		},
		{
			name: "contain empty string in CLIENT_ORIGIN",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"CLIENT_ORIGIN": "http://localhost:8080,,http://localhost:8081",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			out: &config.Config{
				Port:         8080,
				DBHost:       "localhost",
				DBPort:       3306,
				DBName:       "app",
				DBUsername:   "user",
				ClientOrigin: config.ClientOrigin{"http://localhost:8080", "http://localhost:8081"},
				DBUrl:        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				AuthSecret:   "custom-secret",
				SecretIssuer: "custom-issuer",
				AppEnv:       config.ProductionEnv,
				LogLevel:     config.InfoLevel,
			},
		},
		{
			name: "invalid PORT",
			env: map[string]string{
				"PORT":          "invalid",
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			failed: true,
		},
		{
			name: "invalid DB_PORT",
			env: map[string]string{
				"DB_PORT":       "invalid",
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			failed: true,
		},
		{
			name: "invalid FAKE_TIME",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"FAKE_TIME":     "invalid",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			failed: true,
		},
		{
			name: "invalid APP_ENV",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"APP_ENV":       "invalid",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			failed: true,
		},
		{
			name: "invalid LOG_LEVEL",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"LOG_LEVEL":     "invalid",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
			},
			failed: true,
		},
		{
			name: "invalid DEBUG_CORS",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DEBUG_CORS":    "invalid",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			failed: true,
		},
		{
			name: "missing DB_NAME",
			env: map[string]string{
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
				"LOG_LEVEL":     "info",
			},
			failed: true,
		},
		{
			name: "missing DB_URL",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"LOG_LEVEL":     "info",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
			},
			failed: true,
		},
		{
			name: "missing LOG_LEVEL",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"SECRET_ISSUER": "custom-issuer",
			},
			failed: true,
		},
		{
			name: "missing DB_USERNAME",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET":   "custom-secret",
				"LOG_LEVEL":     "info",
				"SECRET_ISSUER": "custom-issuer",
			},
			failed: true,
		},
		{
			name: "missing AUTH_SECRET",
			env: map[string]string{
				"DB_NAME":       "app",
				"DB_USERNAME":   "user",
				"DB_URL":        "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"LOG_LEVEL":     "info",
				"SECRET_ISSUER": "custom-issuer",
			},
			failed: true,
		},
		{
			name: "missing SECRET_ISSUER",
			env: map[string]string{
				"DB_NAME":     "app",
				"DB_USERNAME": "user",
				"DB_URL":      "postgres://postgres:password@localhost:5432/rs?sslmode=disable",
				"AUTH_SECRET": "custom-secret",
				"LOG_LEVEL":   "info",
			},
			failed: true,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			for key, value := range v.env {
				tt.Setenv(key, value)
			}

			cfg, err := config.Get()
			switch {
			case err != nil && !v.failed:
				tt.Fatalf("unexpected error: %+v", err)
			case err == nil && v.failed:
				tt.Fatal("unexpected success")
			case err != nil && v.failed:
				// pass
				tt.Logf("expected error: %+v", err)
				return
			}

			if diff := cmp.Diff(v.out, cfg); diff != "" {
				tt.Errorf("unexpected result:\n%s", diff)
			}
		})
	}
}
