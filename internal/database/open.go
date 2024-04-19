// Package database connection open function.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Open opens a database specified by its database driver name and a driver-specific data source name.
func Open(host string, port uint16, dbname, username, password string) (*pgx.Conn, error) {
	ctx := context.Background()
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, username, password)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return conn, nil
}
