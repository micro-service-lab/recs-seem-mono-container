package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Absence 欠席を表すインターフェース。
type Absence interface {
	// CountAbsences 欠席数を取得する。
	CountAbsences(ctx context.Context) (int64, error)
	// CountAbsencesWithSd SD付きで欠席数を取得する。
	CountAbsencesWithSd(ctx context.Context, sd Sd) (int64, error)
	// CreateAbsence 欠席を作成する。
	CreateAbsence(ctx context.Context, param parameter.CreateAbsenceParam) (entity.Absence, error)
	// CreateAbsenceWithSd SD付きで欠席を作成する。
	CreateAbsenceWithSd(ctx context.Context, sd Sd, param parameter.CreateAbsenceParam) (entity.Absence, error)
	// CreateAbsences 欠席を作成する。
	CreateAbsences(ctx context.Context, params []parameter.CreateAbsenceParam) (int64, error)
	// CreateAbsencesWithSd SD付きで欠席を作成する。
	CreateAbsencesWithSd(ctx context.Context, sd Sd, params []parameter.CreateAbsenceParam) (int64, error)
	// DeleteAbsence 欠席を削除する。
	DeleteAbsence(ctx context.Context, absenceID uuid.UUID) (int64, error)
	// DeleteAbsenceWithSd SD付きで欠席を削除する。
	DeleteAbsenceWithSd(ctx context.Context, sd Sd, absenceID uuid.UUID) (int64, error)
	// PluralDeleteAbsences 欠席を複数削除する。
	PluralDeleteAbsences(ctx context.Context, absenceIDs []uuid.UUID) (int64, error)
	// PluralDeleteAbsencesWithSd SD付きで欠席を複数削除する。
	PluralDeleteAbsencesWithSd(ctx context.Context, sd Sd, absenceIDs []uuid.UUID) (int64, error)
	// FindAbsenceByID 欠席を取得する。
	FindAbsenceByID(ctx context.Context, absenceID uuid.UUID) (entity.Absence, error)
	// FindAbsenceByIDWithSd SD付きで欠席を取得する。
	FindAbsenceByIDWithSd(ctx context.Context, sd Sd, absenceID uuid.UUID) (entity.Absence, error)
	// GetAbsences 欠席を取得する。
	GetAbsences(
		ctx context.Context, order parameter.AbsenceOrderMethod,
		np NumberedPaginationParam, cp CursorPaginationParam, wc WithCountParam,
	) (ListResult[entity.Absence], error)
	// GetAbsencesWithSd SD付きで欠席を取得する。
	GetAbsencesWithSd(
		ctx context.Context, sd Sd, order parameter.AbsenceOrderMethod,
		np NumberedPaginationParam, cp CursorPaginationParam, wc WithCountParam,
	) (ListResult[entity.Absence], error)
	// GetPluralAbsences 欠席を取得する。
	GetPluralAbsences(ctx context.Context, ids []uuid.UUID,
		order parameter.AbsenceOrderMethod, np NumberedPaginationParam) (ListResult[entity.Absence], error)
	// GetPluralAbsencesWithSd SD付きで欠席を取得する。
	GetPluralAbsencesWithSd(ctx context.Context, sd Sd, ids []uuid.UUID,
		order parameter.AbsenceOrderMethod,
		np NumberedPaginationParam) (ListResult[entity.Absence], error)
}
