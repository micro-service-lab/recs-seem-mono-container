package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateAttendStatusParam 出席ステータス作成のパラメータ。
type CreateAttendStatusParam struct {
	Name string
	Key  string
}

// UpdateAttendStatusParams 出席ステータス更新のパラメータ。
type UpdateAttendStatusParams struct {
	Name string
	Key  string
}

// UpdateAttendStatusByKeyParams 出席ステータス更新のパラメータ。
type UpdateAttendStatusByKeyParams struct {
	Name string
}

// WhereAttendStatusParam 出席ステータス検索のパラメータ。
type WhereAttendStatusParam struct {
	WhereLikeName bool
	SearchName    string
}

// AttendStatusOrderMethod 出席ステータスの並び替え方法。
type AttendStatusOrderMethod string

// IsMatch は文字列が一致するかを返す。
func (a AttendStatusOrderMethod) IsMatch(s string) bool {
	return a == AttendStatusOrderMethod(s)
}

const (
	// AttendStatusOrderMethodName は名前順。
	AttendStatusOrderMethodName AttendStatusOrderMethod = "name"
	// AttendStatusOrderMethodReverseName は名前逆順。
	AttendStatusOrderMethodReverseName AttendStatusOrderMethod = "r_name"
)

// AttendStatus 出席ステータスを表すインターフェース。
type AttendStatus interface {
	// CountAttendStatuses 出席ステータス数を取得する。
	CountAttendStatuses(ctx context.Context, where WhereAttendStatusParam) (int64, error)
	// CountAttendStatusesWithSd SD付きで出席ステータス数を取得する。
	CountAttendStatusesWithSd(ctx context.Context, sd Sd, where WhereAttendStatusParam) (int64, error)
	// CreateAttendStatus 出席ステータスを作成する。
	CreateAttendStatus(ctx context.Context, param CreateAttendStatusParam) (entity.AttendStatus, error)
	// CreateAttendStatusWithSd SD付きで出席ステータスを作成する。
	CreateAttendStatusWithSd(ctx context.Context, sd Sd, param CreateAttendStatusParam) (entity.AttendStatus, error)
	// CreateAttendStatuses 出席ステータスを作成する。
	CreateAttendStatuses(ctx context.Context, params []CreateAttendStatusParam) (int64, error)
	// CreateAttendStatusesWithSd SD付きで出席ステータスを作成する。
	CreateAttendStatusesWithSd(ctx context.Context, sd Sd, params []CreateAttendStatusParam) (int64, error)
	// DeleteAttendStatus 出席ステータスを削除する。
	DeleteAttendStatus(ctx context.Context, attendStatusID uuid.UUID) error
	// DeleteAttendStatusWithSd SD付きで出席ステータスを削除する。
	DeleteAttendStatusWithSd(ctx context.Context, sd Sd, attendStatusID uuid.UUID) error
	// DeleteAttendStatusByKey 出席ステータスを削除する。
	DeleteAttendStatusByKey(ctx context.Context, key string) error
	// DeleteAttendStatusByKeyWithSd SD付きで出席ステータスを削除する。
	DeleteAttendStatusByKeyWithSd(ctx context.Context, sd Sd, key string) error
	// FindAttendStatusByID 出席ステータスを取得する。
	FindAttendStatusByID(ctx context.Context, attendStatusID uuid.UUID) (entity.AttendStatus, error)
	// FindAttendStatusByIDWithSd SD付きで出席ステータスを取得する。
	FindAttendStatusByIDWithSd(ctx context.Context, sd Sd, attendStatusID uuid.UUID) (entity.AttendStatus, error)
	// FindAttendStatusByKey 出席ステータスを取得する。
	FindAttendStatusByKey(ctx context.Context, key string) (entity.AttendStatus, error)
	// FindAttendStatusByKeyWithSd SD付きで出席ステータスを取得する。
	FindAttendStatusByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.AttendStatus, error)
	// GetAttendStatuses 出席ステータスを取得する。
	GetAttendStatuses(
		ctx context.Context,
		where WhereAttendStatusParam,
		order AttendStatusOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttendStatus], error)
	// GetAttendStatusesWithSd SD付きで出席ステータスを取得する。
	GetAttendStatusesWithSd(
		ctx context.Context,
		sd Sd,
		where WhereAttendStatusParam,
		order AttendStatusOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.AttendStatus], error)
	// UpdateAttendStatus 出席ステータスを更新する。
	UpdateAttendStatus(
		ctx context.Context, attendStatusID uuid.UUID, param UpdateAttendStatusParams) (entity.AttendStatus, error)
	// UpdateAttendStatusWithSd SD付きで出席ステータスを更新する。
	UpdateAttendStatusWithSd(
		ctx context.Context, sd Sd, attendStatusID uuid.UUID, param UpdateAttendStatusParams) (entity.AttendStatus, error)
	// UpdateAttendStatusByKey 出席ステータスを更新する。
	UpdateAttendStatusByKey(
		ctx context.Context, key string, param UpdateAttendStatusByKeyParams) (entity.AttendStatus, error)
	// UpdateAttendStatusByKeyWithSd SD付きで出席ステータスを更新する。
	UpdateAttendStatusByKeyWithSd(
		ctx context.Context, sd Sd, key string, param UpdateAttendStatusByKeyParams) (entity.AttendStatus, error)
}
