package pgadapter_test

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store/pgadapter"
)

// SeriesTest is a type for table driven tests.
type SeriesTest struct {
	Name string
	Test func(t *testing.T)
}

// NewDummyPgAdapter creates a new PgAdapter.
func NewDummyPgAdapter(t *testing.T) *pgadapter.PgAdapter {
	t.Helper()

	ctx := context.Background()

	var initScripts []string

	wd, err := os.Getwd()
	require.NoError(t, err)

	migrationPath := filepath.Join(wd, "../../../", "db", "migrations")
	files, err := os.ReadDir(migrationPath)
	require.NoError(t, err)

	for _, file := range files {
		// .up.sql files are executed in ascending order
		if strings.HasSuffix(file.Name(), ".up.sql") {
			initScripts = append(initScripts, filepath.Join(migrationPath, file.Name()))
		}
	}
	// file name order
	sort.Strings(initScripts)

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithInitScripts(initScripts...),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			require.NoError(t, err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	adapter, err := pgadapter.NewPgAdapterFromConnStr(ctx, connStr, nil)
	if err != nil {
		require.NoError(t, err)
	}

	return adapter
}

var validNp = store.NumberedPaginationParam{
	Valid:  true,
	Offset: entity.Int{Int64: 0, Valid: true},
	Limit:  entity.Int{Int64: 10, Valid: true},
}

// var invalidNp = store.NumberedPaginationParam{}

// var validCp = store.CursorPaginationParam{
// 	Valid:  true,
// 	Cursor: "",
// 	Limit:  entity.Int{Int64: 10, Valid: true},
// }

var invalidCp = store.CursorPaginationParam{}

var validWc = store.WithCountParam{
	Valid: true,
}

// var invalidWc = store.WithCountParam{}
