// Package store 永続化関連のアダプタを提供するパッケージ
package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ErrNotFoundDescriptor ディスクリプタが見つからないエラー。
var ErrNotFoundDescriptor = fmt.Errorf("not found descriptor")

// Sd Storeディスクリプタ。
type Sd uuid.UUID

// Store 永続化関連のアダプタを提供するインターフェース。
type Store interface {
	// WithTx トランザクションを開始する。
	Begin(context.Context) (Sd, error)
	// Commit トランザクションをコミットする。
	Commit(context.Context, Sd) error
	// Rollback トランザクションをロールバックする。
	Rollback(context.Context, Sd) error
	// Cleanup ストアをクリーンアップする。
	Cleanup(context.Context) error

	Absence
	AttendStatus
}
