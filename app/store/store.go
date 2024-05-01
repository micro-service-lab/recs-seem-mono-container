// Package store 永続化関連のアダプタを提供するパッケージ
package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
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

// ListResult リストの結果を表す型。
type ListResult[T entity.Entity] struct {
	// Data エンティティのスライス。
	Data []T `json:"data"`
	// CursorPagination ページネーション情報。
	CursorPagination CursorPaginationAttribute `json:"cursor_pagination"`
	// WithCount カウント情報。
	WithCount WithCountAttribute `json:"with_count"`
}

type (
	// EntityConvertFunc エンティティ変換関数。
	EntityConvertFunc[T any, U entity.Entity] func(T) (U, error)
	// RunGetCountFunc カウントを取得する関数。
	RunGetCountFunc func() (int64, error)
	// RunQueryFunc クエリを実行する関数。
	RunQueryFunc[T any] func(orderMethod string) ([]T, error)
	// RunQueryWithCursorParamsFunc カーソルパラメータを持つクエリを実行する関数。
	RunQueryWithCursorParamsFunc[T any] func(subCursor string, orderMethod string, limit int32, cursorDir string, cursor int32, subCursorValue any) ([]T, error)
	// RunQueryWithNumberedParamsFunc ページネーションパラメータを持つクエリを実行する関数。
	RunQueryWithNumberedParamsFunc[T any] func(orderMethod string, limit int32, offset int32) ([]T, error)
	// CursorIDAndValueSelector カーソルIDと値を選択する関数。
	CursorIDAndValueSelector[T any] func(subCursor string, e T) (entity.Int, any)
)

// RunListQuery クエリを実行する。
func RunListQuery[T any, U entity.Entity](
	_ context.Context,
	order parameter.OrderMethod,
	np NumberedPaginationParam,
	cp CursorPaginationParam,
	wc WithCountParam,
	eConvFunc EntityConvertFunc[T, U],
	runCFunc RunGetCountFunc,
	runQFunc RunQueryFunc[T],
	runQCPFunc RunQueryWithCursorParamsFunc[T],
	runQNPFunc RunQueryWithNumberedParamsFunc[T],
	selector CursorIDAndValueSelector[T],
) (ListResult[U], error) {
	var withCount int64
	var err error
	if wc.Valid {
		withCount, err = runCFunc()
		if err != nil {
			return ListResult[U]{}, fmt.Errorf("failed to get count: %w", err)
		}
	}
	wcAtr := WithCountAttribute{
		Count: withCount,
		Valid: wc.Valid,
	}

	if np.Valid {
		e, err := runQNPFunc(order.GetStringValue(), int32(np.Limit.Int64), int32(np.Offset.Int64))
		if err != nil {
			return ListResult[U]{}, fmt.Errorf("failed to get entities with numbered pagination: %w", err)
		}
		entities := make([]U, len(e))
		for i, v := range e {
			entities[i], err = eConvFunc(v)
			if err != nil {
				return ListResult[U]{}, fmt.Errorf("failed to convert entity: %w", err)
			}
		}
		return ListResult[U]{Data: entities, WithCount: wcAtr}, nil
	} else if cp.Valid {
		e, pi, err := GetCursorData(cp.Cursor, order, int32(cp.Limit.Int64), runQCPFunc, runQNPFunc, selector)
		if err != nil {
			return ListResult[U]{}, fmt.Errorf("failed to get entities with cursor pagination: %w", err)
		}
		entities := make([]U, len(e))
		for i, v := range e {
			entities[i], err = eConvFunc(v)
			if err != nil {
				return ListResult[U]{}, fmt.Errorf("failed to convert entity: %w", err)
			}
		}
		return ListResult[U]{Data: entities, CursorPagination: pi, WithCount: wcAtr}, nil
	}
	e, err := runQFunc(order.GetStringValue())
	if err != nil {
		return ListResult[U]{}, fmt.Errorf("failed to get entities: %w", err)
	}
	entities := make([]U, len(e))
	for i, v := range e {
		entities[i], err = eConvFunc(v)
		if err != nil {
			return ListResult[U]{}, fmt.Errorf("failed to convert entity: %w", err)
		}
	}
	return ListResult[U]{Data: entities, WithCount: wcAtr}, nil
}
