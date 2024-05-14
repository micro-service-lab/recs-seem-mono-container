package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// WorkPosition ワークポジションを表すインターフェース。
type WorkPosition interface {
	// CountWorkPositions ワークポジション数を取得する。
	CountWorkPositions(ctx context.Context, where parameter.WhereWorkPositionParam) (int64, error)
	// CountWorkPositionsWithSd SD付きでワークポジション数を取得する。
	CountWorkPositionsWithSd(ctx context.Context, sd Sd, where parameter.WhereWorkPositionParam) (int64, error)
	// CreateWorkPosition ワークポジションを作成する。
	CreateWorkPosition(ctx context.Context, param parameter.CreateWorkPositionParam) (entity.WorkPosition, error)
	// CreateWorkPositionWithSd SD付きでワークポジションを作成する。
	CreateWorkPositionWithSd(
		ctx context.Context, sd Sd, param parameter.CreateWorkPositionParam) (entity.WorkPosition, error)
	// CreateWorkPositions ワークポジションを作成する。
	CreateWorkPositions(ctx context.Context, params []parameter.CreateWorkPositionParam) (int64, error)
	// CreateWorkPositionsWithSd SD付きでワークポジションを作成する。
	CreateWorkPositionsWithSd(ctx context.Context, sd Sd, params []parameter.CreateWorkPositionParam) (int64, error)
	// DeleteWorkPosition ワークポジションを削除する。
	DeleteWorkPosition(ctx context.Context, roleID uuid.UUID) (int64, error)
	// DeleteWorkPositionWithSd SD付きでワークポジションを削除する。
	DeleteWorkPositionWithSd(ctx context.Context, sd Sd, roleID uuid.UUID) (int64, error)
	// PluralDeleteWorkPositions ワークポジションを複数削除する。
	PluralDeleteWorkPositions(ctx context.Context, roleIDs []uuid.UUID) (int64, error)
	// PluralDeleteWorkPositionsWithSd SD付きでワークポジションを複数削除する。
	PluralDeleteWorkPositionsWithSd(ctx context.Context, sd Sd, roleIDs []uuid.UUID) (int64, error)
	// FindWorkPositionByID ワークポジションを取得する。
	FindWorkPositionByID(ctx context.Context, roleID uuid.UUID) (entity.WorkPosition, error)
	// FindWorkPositionByIDWithSd SD付きでワークポジションを取得する。
	FindWorkPositionByIDWithSd(ctx context.Context, sd Sd, roleID uuid.UUID) (entity.WorkPosition, error)
	// GetWorkPositions ワークポジションを取得する。
	GetWorkPositions(
		ctx context.Context,
		where parameter.WhereWorkPositionParam,
		order parameter.WorkPositionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.WorkPosition], error)
	// GetWorkPositionsWithSd SD付きでワークポジションを取得する。
	GetWorkPositionsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereWorkPositionParam,
		order parameter.WorkPositionOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.WorkPosition], error)
	// GetPluralWorkPositions ワークポジションを取得する。
	GetPluralWorkPositions(
		ctx context.Context,
		workPositionIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.WorkPosition], error)
	// GetPluralWorkPositionsWithSd SD付きでワークポジションを取得する。
	GetPluralWorkPositionsWithSd(
		ctx context.Context,
		sd Sd,
		workPositionIDs []uuid.UUID,
		np NumberedPaginationParam,
	) (ListResult[entity.WorkPosition], error)
	// UpdateWorkPosition ワークポジションを更新する。
	UpdateWorkPosition(
		ctx context.Context,
		roleID uuid.UUID,
		param parameter.UpdateWorkPositionParams,
	) (entity.WorkPosition, error)
	// UpdateWorkPositionWithSd SD付きでワークポジションを更新する。
	UpdateWorkPositionWithSd(
		ctx context.Context, sd Sd, roleID uuid.UUID,
		param parameter.UpdateWorkPositionParams) (entity.WorkPosition, error)
}
