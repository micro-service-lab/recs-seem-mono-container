package batch

import (
	"context"
	"fmt"

	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
)

// InitAttendanceTypes is a batch to initialize attendance types.
type InitAttendanceTypes struct {
	Manager *service.ManagerInterface
}

// Run initializes attendance types.
func (c *InitAttendanceTypes) Run(ctx context.Context) error {
	var as []parameter.CreateAttendanceTypeParam
	for _, a := range service.AttendanceTypes {
		as = append(as, parameter.CreateAttendanceTypeParam{
			Name:  a.Name,
			Key:   a.Key,
			Color: a.Color,
		})
	}
	_, err := (*c.Manager).CreateAttendanceTypes(ctx, as)
	if err != nil {
		return fmt.Errorf("failed to create attendance types: %w", err)
	}
	return nil
}
