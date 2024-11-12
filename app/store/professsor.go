package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Professor 教授を表すインターフェース。
type Professor interface {
	// CountProfessors 教授数を取得する。
	CountProfessors(ctx context.Context, where parameter.WhereProfessorParam) (int64, error)
	// CountProfessorsWithSd SD付きで教授数を取得する。
	CountProfessorsWithSd(ctx context.Context, sd Sd, where parameter.WhereProfessorParam) (int64, error)
	// CreateProfessor 教授を作成する。
	CreateProfessor(ctx context.Context, param parameter.CreateProfessorParam) (entity.Professor, error)
	// CreateProfessorWithSd SD付きで教授を作成する。
	CreateProfessorWithSd(
		ctx context.Context, sd Sd, param parameter.CreateProfessorParam) (entity.Professor, error)
	// CreateProfessors 教授を作成する。
	CreateProfessors(ctx context.Context, params []parameter.CreateProfessorParam) (int64, error)
	// CreateProfessorsWithSd SD付きで教授を作成する。
	CreateProfessorsWithSd(ctx context.Context, sd Sd, params []parameter.CreateProfessorParam) (int64, error)
	// DeleteProfessor 教授を削除する。
	DeleteProfessor(ctx context.Context, professorID uuid.UUID) (int64, error)
	// DeleteProfessorWithSd SD付きで教授を削除する。
	DeleteProfessorWithSd(ctx context.Context, sd Sd, professorID uuid.UUID) (int64, error)
	// PluralDeleteProfessors 教授を複数削除する。
	PluralDeleteProfessors(ctx context.Context, professorIDs []uuid.UUID) (int64, error)
	// PluralDeleteProfessorsWithSd SD付きで教授を複数削除する。
	PluralDeleteProfessorsWithSd(ctx context.Context, sd Sd, professorIDs []uuid.UUID) (int64, error)
	// FindProfessorByID 教授を取得する。
	FindProfessorByID(ctx context.Context, professorID uuid.UUID) (entity.Professor, error)
	// FindProfessorByIDWithSd SD付きで教授を取得する。
	FindProfessorByIDWithSd(ctx context.Context, sd Sd, professorID uuid.UUID) (entity.Professor, error)
	// FindProfessorWithMember 教授を取得する。
	FindProfessorWithMember(ctx context.Context, professorID uuid.UUID) (entity.ProfessorWithMember, error)
	// FindProfessorWithMemberWithSd SD付きで教授を取得する。
	FindProfessorWithMemberWithSd(
		ctx context.Context, sd Sd, professorID uuid.UUID) (entity.ProfessorWithMember, error)
	// GetProfessors 教授を取得する。
	GetProfessors(
		ctx context.Context,
		where parameter.WhereProfessorParam,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Professor], error)
	// GetProfessorsWithSd SD付きで教授を取得する。
	GetProfessorsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereProfessorParam,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Professor], error)
	// GetPluralProfessors 教授を取得する。
	GetPluralProfessors(
		ctx context.Context,
		professorIDs []uuid.UUID,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Professor], error)
	// GetPluralProfessorsWithSd SD付きで教授を取得する。
	GetPluralProfessorsWithSd(
		ctx context.Context,
		sd Sd,
		professorIDs []uuid.UUID,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Professor], error)
	// GetProfessorsWithMember 教授を取得する。
	GetProfessorsWithMember(
		ctx context.Context,
		where parameter.WhereProfessorParam,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ProfessorWithMember], error)
	// GetProfessorsWithMemberWithSd SD付きで教授を取得する。
	GetProfessorsWithMemberWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereProfessorParam,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.ProfessorWithMember], error)
	// GetPluralProfessorsWithMember 教授を取得する。
	GetPluralProfessorsWithMember(
		ctx context.Context,
		professorIDs []uuid.UUID,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ProfessorWithMember], error)
	// GetPluralProfessorsWithMemberWithSd SD付きで教授を取得する。
	GetPluralProfessorsWithMemberWithSd(
		ctx context.Context,
		sd Sd,
		professorIDs []uuid.UUID,
		order parameter.ProfessorOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.ProfessorWithMember], error)
}
