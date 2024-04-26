package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// ManageAbsence 欠席管理サービス。
type ManageAbsence struct {
	store store.Store
}

// CreateAbsence 欠席を作成する。
func (s *ManageAbsence) CreateAbsence(
	ctx context.Context,
	attendanceID uuid.UUID,
) (entity.Absence, error) {
	p := store.CreateAbsenceParam{
		AttendanceID: attendanceID,
	}
	e, err := s.store.CreateAbsence(ctx, p)
	if err != nil {
		return entity.Absence{}, fmt.Errorf("failed to create absence: %w", err)
	}
	return e, nil
}
