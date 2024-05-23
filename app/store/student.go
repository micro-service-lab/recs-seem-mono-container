package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
)

// Student 生徒を表すインターフェース。
type Student interface {
	// CountStudents 生徒数を取得する。
	CountStudents(ctx context.Context, where parameter.WhereStudentParam) (int64, error)
	// CountStudentsWithSd SD付きで生徒数を取得する。
	CountStudentsWithSd(ctx context.Context, sd Sd, where parameter.WhereStudentParam) (int64, error)
	// CreateStudent 生徒を作成する。
	CreateStudent(ctx context.Context, param parameter.CreateStudentParam) (entity.Student, error)
	// CreateStudentWithSd SD付きで生徒を作成する。
	CreateStudentWithSd(
		ctx context.Context, sd Sd, param parameter.CreateStudentParam) (entity.Student, error)
	// CreateStudents 生徒を作成する。
	CreateStudents(ctx context.Context, params []parameter.CreateStudentParam) (int64, error)
	// CreateStudentsWithSd SD付きで生徒を作成する。
	CreateStudentsWithSd(ctx context.Context, sd Sd, params []parameter.CreateStudentParam) (int64, error)
	// DeleteStudent 生徒を削除する。
	DeleteStudent(ctx context.Context, studentID uuid.UUID) (int64, error)
	// DeleteStudentWithSd SD付きで生徒を削除する。
	DeleteStudentWithSd(ctx context.Context, sd Sd, studentID uuid.UUID) (int64, error)
	// PluralDeleteStudents 生徒を複数削除する。
	PluralDeleteStudents(ctx context.Context, studentIDs []uuid.UUID) (int64, error)
	// PluralDeleteStudentsWithSd SD付きで生徒を複数削除する。
	PluralDeleteStudentsWithSd(ctx context.Context, sd Sd, studentIDs []uuid.UUID) (int64, error)
	// FindStudentByID 生徒を取得する。
	FindStudentByID(ctx context.Context, studentID uuid.UUID) (entity.Student, error)
	// FindStudentByIDWithSd SD付きで生徒を取得する。
	FindStudentByIDWithSd(ctx context.Context, sd Sd, studentID uuid.UUID) (entity.Student, error)
	// FindStudentWithMember 生徒を取得する。
	FindStudentWithMember(ctx context.Context, studentID uuid.UUID) (entity.StudentWithMember, error)
	// FindStudentWithMemberWithSd SD付きで生徒を取得する。
	FindStudentWithMemberWithSd(
		ctx context.Context, sd Sd, studentID uuid.UUID) (entity.StudentWithMember, error)
	// GetStudents 生徒を取得する。
	GetStudents(
		ctx context.Context,
		where parameter.WhereStudentParam,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Student], error)
	// GetStudentsWithSd SD付きで生徒を取得する。
	GetStudentsWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereStudentParam,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.Student], error)
	// GetPluralStudents 生徒を取得する。
	GetPluralStudents(
		ctx context.Context,
		studentIDs []uuid.UUID,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Student], error)
	// GetPluralStudentsWithSd SD付きで生徒を取得する。
	GetPluralStudentsWithSd(
		ctx context.Context,
		sd Sd,
		studentIDs []uuid.UUID,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.Student], error)
	// GetStudentsWithMember 生徒を取得する。
	GetStudentsWithMember(
		ctx context.Context,
		where parameter.WhereStudentParam,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.StudentWithMember], error)
	// GetStudentsWithMemberWithSd SD付きで生徒を取得する。
	GetStudentsWithMemberWithSd(
		ctx context.Context,
		sd Sd,
		where parameter.WhereStudentParam,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
		cp CursorPaginationParam,
		wc WithCountParam,
	) (ListResult[entity.StudentWithMember], error)
	// GetPluralStudentsWithMember 生徒を取得する。
	GetPluralStudentsWithMember(
		ctx context.Context,
		studentIDs []uuid.UUID,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.StudentWithMember], error)
	// GetPluralStudentsWithMemberWithSd SD付きで生徒を取得する。
	GetPluralStudentsWithMemberWithSd(
		ctx context.Context,
		sd Sd,
		studentIDs []uuid.UUID,
		order parameter.StudentOrderMethod,
		np NumberedPaginationParam,
	) (ListResult[entity.StudentWithMember], error)
}
