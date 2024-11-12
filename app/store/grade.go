package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Grade 学年を表すインターフェース。
type Grade interface {
	// CountGrades 学年数を取得する。
	CountGrades(ctx context.Context, where parameter.WhereGradeParam) (int64, error)
	// CountGradesWithSd SD付きで学年数を取得する。
	CountGradesWithSd(ctx context.Context, sd Sd, where parameter.WhereGradeParam) (int64, error)
	// CreateGrade 学年を作成する。
	CreateGrade(ctx context.Context, param parameter.CreateGradeParam) (entity.Grade, error)
	// CreateGradeWithSd SD付きで学年を作成する。
	CreateGradeWithSd(
		ctx context.Context, sd Sd, param parameter.CreateGradeParam) (entity.Grade, error)
	// CreateGrades 学年を作成する。
	CreateGrades(ctx context.Context, params []parameter.CreateGradeParam) (int64, error)
	// CreateGradesWithSd SD付きで学年を作成する。
	CreateGradesWithSd(ctx context.Context, sd Sd, params []parameter.CreateGradeParam) (int64, error)
	// DeleteGrade 学年を削除する。
	DeleteGrade(ctx context.Context, gradeID uuid.UUID) (int64, error)
	// DeleteGradeWithSd SD付きで学年を削除する。
	DeleteGradeWithSd(ctx context.Context, sd Sd, gradeID uuid.UUID) (int64, error)
	// PluralDeleteGrades 学年を複数削除する。
	PluralDeleteGrades(ctx context.Context, gradeIDs []uuid.UUID) (int64, error)
	// PluralDeleteGradesWithSd SD付きで学年を複数削除する。
	PluralDeleteGradesWithSd(ctx context.Context, sd Sd, gradeIDs []uuid.UUID) (int64, error)
	// FindGradeByID 学年を取得する。
	FindGradeByID(ctx context.Context, gradeID uuid.UUID) (entity.Grade, error)
	// FindGradeByIDWithSd SD付きで学年を取得する。
	FindGradeByIDWithSd(ctx context.Context, sd Sd, gradeID uuid.UUID) (entity.Grade, error)
	// FindGradeByKey 学年を取得する。
	FindGradeByKey(ctx context.Context, key string) (entity.Grade, error)
	// FindGradeByKeyWithSd SD付きで学年を取得する。
	FindGradeByKeyWithSd(ctx context.Context, sd Sd, key string) (entity.Grade, error)
	// FindGradeWithOrganization 学年を取得する。
	FindGradeWithOrganization(ctx context.Context, gradeID uuid.UUID) (entity.GradeWithOrganization, error)
	// FindGradeWithOrganizationWithSd SD付きで学年を取得する。
	FindGradeWithOrganizationWithSd(
		ctx context.Context, sd Sd, gradeID uuid.UUID) (entity.GradeWithOrganization, error)
	// GetGrades 学年を取得する。
	GetGrades(
		ctx context.Context,
		where parameter.WhereGradeParam,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Grade], error)
	// GetGradesWithSd SD付きで学年を取得する。
	GetGradesWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereGradeParam,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Grade], error)
	// GetPluralGrades 学年を取得する。
	GetPluralGrades(
		ctx context.Context,
		gradeIDs []uuid.UUID,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Grade], error)
	// GetPluralGradesWithSd SD付きで学年を取得する。
	GetPluralGradesWithSd(
		ctx context.Context,
		sd Sd,
		gradeIDs []uuid.UUID,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Grade], error)
	// GetGradesWithOrganization 学年を取得する。
	GetGradesWithOrganization(
		ctx context.Context,
		where parameter.WhereGradeParam,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.GradeWithOrganization], error)
	// GetGradesWithOrganizationWithSd SD付きで学年を取得する。
	GetGradesWithOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereGradeParam,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.GradeWithOrganization], error)
	// GetPluralGradesWithOrganization 学年を取得する。
	GetPluralGradesWithOrganization(
		ctx context.Context,
		gradeIDs []uuid.UUID,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.GradeWithOrganization], error)
	// GetPluralGradesWithOrganizationWithSd SD付きで学年を取得する。
	GetPluralGradesWithOrganizationWithSd(
		ctx context.Context,
		sd Sd,
		gradeIDs []uuid.UUID,
		order parameter.GradeOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.GradeWithOrganization], error)
}
