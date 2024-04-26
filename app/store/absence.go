package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
)

// CreateAbsenceParam 欠席作成のパラメータ。
type CreateAbsenceParam struct {
	AttendanceID uuid.UUID
}

// Absence 欠席を表すインターフェース。
type Absence interface {
	// CountAbsences 欠席数を取得する。
	CountAbsences(ctx context.Context) (int64, error)
	// CountAbsencesWithSd SD付きで欠席数を取得する。
	CountAbsencesWithSd(ctx context.Context, sd Sd) (int64, error)
	// CreateAbsence 欠席を作成する。
	CreateAbsence(ctx context.Context, param CreateAbsenceParam) (entity.Absence, error)
	// CreateAbsenceWithSd SD付きで欠席を作成する。
	CreateAbsenceWithSd(ctx context.Context, sd Sd, param CreateAbsenceParam) (entity.Absence, error)
	// CreateAbsences 欠席を作成する。
	CreateAbsences(ctx context.Context, params []CreateAbsenceParam) (int64, error)
	// CreateAbsencesWithSd SD付きで欠席を作成する。
	CreateAbsencesWithSd(ctx context.Context, sd Sd, params []CreateAbsenceParam) (int64, error)
	// DeleteAbsence 欠席を削除する。
	DeleteAbsence(ctx context.Context, absenceID uuid.UUID) error
	// DeleteAbsenceWithSd SD付きで欠席を削除する。
	DeleteAbsenceWithSd(ctx context.Context, sd Sd, absenceID uuid.UUID) error
	// FindAbsenceByID 欠席を取得する。
	FindAbsenceByID(ctx context.Context, absenceID uuid.UUID) (entity.Absence, error)
	// FindAbsenceByIDWithSd SD付きで欠席を取得する。
	FindAbsenceByIDWithSd(ctx context.Context, sd Sd, absenceID uuid.UUID) (entity.Absence, error)
	// GetAbsences 欠席を取得する。
	GetAbsences(
		ctx context.Context, np NumberedPaginationParam, cp CursorPaginationParam, wc WithCountParam,
	) (ListResult[entity.Absence], error)
	// GetAbsencesWithSd SD付きで欠席を取得する。
	GetAbsencesWithSd(
		ctx context.Context, sd Sd, np NumberedPaginationParam, cp CursorPaginationParam, wc WithCountParam,
	) (ListResult[entity.Absence], error)
}
