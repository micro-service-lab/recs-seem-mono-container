package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// GradeKey 年次キー。
type GradeKey string

const (
	// GradeKeyBachelor1 1年生。
	GradeKeyBachelor1 GradeKey = "bachelor1"
	// GradeKeyBachelor2 2年生。
	GradeKeyBachelor2 GradeKey = "bachelor2"
	// GradeKeyBachelor3 3年生。
	GradeKeyBachelor3 GradeKey = "bachelor3"
	// GradeKeyBachelor4 4年生。
	GradeKeyBachelor4 GradeKey = "bachelor4"
	// GradeKeyMaster1 修士1年生。
	GradeKeyMaster1 GradeKey = "master1"
	// GradeKeyMaster2 修士2年生。
	GradeKeyMaster2 GradeKey = "master2"
	// GradeKeyDoctor 博士。
	GradeKeyDoctor GradeKey = "doctor"
	// GradeKeyProfessor 教授。
	GradeKeyProfessor GradeKey = "professor"
)

// Grade 年次。
type Grade struct {
	Key         string
	Name        string
	Description string
	Color       string
}

// Grades 年次一覧。
var Grades = []Grade{
	{
		Key:         string(GradeKeyBachelor1),
		Name:        "B1",
		Description: "学部1回生",
		Color:       "#FF0000",
	},
	{
		Key:         string(GradeKeyBachelor2),
		Name:        "B2",
		Description: "学部2回生",
		Color:       "#00FF00",
	},
	{
		Key:         string(GradeKeyBachelor3),
		Name:        "B3",
		Description: "学部3回生",
		Color:       "#0000FF",
	},
	{
		Key:         string(GradeKeyBachelor4),
		Name:        "B4",
		Description: "学部4回生",
		Color:       "#FFFF00",
	},
	{
		Key:         string(GradeKeyMaster1),
		Name:        "M1",
		Description: "修士1回生",
		Color:       "#FF00FF",
	},
	{
		Key:         string(GradeKeyMaster2),
		Name:        "M2",
		Description: "修士2回生",
		Color:       "#00FFFF",
	},
	{
		Key:         string(GradeKeyDoctor),
		Name:        "ドクター",
		Description: "博士",
		Color:       "#000000",
	},
	{
		Key:         string(GradeKeyProfessor),
		Name:        "教授",
		Description: "教授",
		Color:       "#FFFFFF",
	},
}

// ManageGrade 年次管理サービス。
type ManageGrade struct {
	DB store.Store
}

// CreateGrade 年次を作成する。
func (m *ManageGrade) CreateGrade(
	ctx context.Context,
	name, key, description string,
	categoryID uuid.UUID,
) (entity.Grade, error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	p := parameter.CreateGradeParam{
		// Name:            name,
		// Key:             key,
		// Description:     description,
		// GradeCategoryID: categoryID,
	}
	e, err := m.DB.CreateGrade(ctx, p)
	if err != nil {
		return entity.Grade{}, fmt.Errorf("failed to create policy: %w", err)
	}
	return e, nil
}

// CreateGrades 年次を複数作成する。
func (m *ManageGrade) CreateGrades(
	ctx context.Context, ps []parameter.CreateGradeParam,
) (int64, error) {
	es, err := m.DB.CreateGrades(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create policy: %w", err)
	}
	return es, nil
}

// DeleteGrade 年次を削除する。
func (m *ManageGrade) DeleteGrade(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeleteGrade(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy: %w", err)
	}
	return c, nil
}

// PluralDeleteGrades 年次を複数削除する。
func (m *ManageGrade) PluralDeleteGrades(
	ctx context.Context, ids []uuid.UUID,
) (int64, error) {
	c, err := m.DB.PluralDeleteGrades(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete policy: %w", err)
	}
	return c, nil
}
