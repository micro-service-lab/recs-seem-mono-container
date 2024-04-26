// Package service provides a application service.
package service

import (
	"context"

	"github.com/google/uuid"
)

// AbsenceManager is a interface for absence service.
type AbsenceManager interface {
	CreateAbsence(ctx context.Context, attendanceID uuid.UUID) error
}
