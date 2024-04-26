// Package store 永続化関連のアダプタを提供するパッケージ
package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
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

	Absence
	AttendStatus
}

// ListResult リストの結果を表す型。
type ListResult[T entity.Entity] struct {
	// Data エンティティのスライス。
	Data []T `json:"data"`
	// CursorPagination ページネーション情報。
	CursorPagination CursorPaginationAttribute `json:"cursor_pagination"`
	// WithCount カウント情報。
	WithCount WithCountAttribute `json:"with_count"`
}
