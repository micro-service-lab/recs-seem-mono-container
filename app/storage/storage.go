// Package storage provides a simple interface to store and retrieve data
package storage

import (
	"context"
	"io"
)

// Storage is an interface to store and retrieve data.
type Storage interface {
	// PutObject stores data and returns the key.
	PutObject(ctx context.Context, reader io.Reader, key, contentType string) (string, error)
	// UploadObject stores data and returns the key.
	UploadObject(
		ctx context.Context, reader io.Reader, key, contentType string,
	) (string, error)
	// GetObject retrieves data by key.
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	// DownloadObject retrieves data by key.
	DownloadObject(ctx context.Context, writer io.WriterAt, key string) error
	// DeleteObjects deletes data by key.
	DeleteObjects(ctx context.Context, keys []string) error
	// ExistsObject checks if data exists by key.
	ExistsObject(ctx context.Context, key string) (bool, error)
	// GetURLFromKey returns the URL from the key.
	GetURLFromKey(ctx context.Context, key string) (string, error)
	// GetKeyFromURL returns the key from the URL.
	GetKeyFromURL(ctx context.Context, url string) (string, error)
}
